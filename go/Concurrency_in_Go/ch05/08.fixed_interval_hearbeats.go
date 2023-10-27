package main

import (
	"testing"
	"time"
)

func DoWork(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		// 인수로 들어온 pulse 간격으로 pulse에서 데이터 보내기
		pulse := time.Tick(pulseInterval)
	numLoop:
		for _, n := range nums {
			for {
				select {
				// done 채널에 데이터가 들어오면 고루틴 종료
				case <-done:
					return
				// pulse 간격에 따라 신호를 받은 경우
				case <-pulse:
					select {
					// 하트비트에 신호 보내기
					case heartbeat <- struct{}{}:
					default:
					}
				// 인수로 들어온 숫자 슬라이스 intStream 채널에 넣기
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeat, intStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	// 1초마다 하트비트 채널, 작업 채널로 값 받기
	heartbeat, results := DoWork(done, timeout/2, intSlice...)

	// 하트비트 채널에 신호가 들어오면
	<-heartbeat

	i := 0
	for {
		select {
		// 값이 들어오면
		case r, ok := <-results:
			// 채널에 값이 없고 닫힌 상태라면
			if ok == false {
				// 테스트 종료
				return
				// 보낸 숫자 슬라이스와 DoWork 채널에서 받은 데이터 비교
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat:
		// 작업 채널에서 값이 안넘어오고 2초가 지난 경우 타임아웃 에러
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
