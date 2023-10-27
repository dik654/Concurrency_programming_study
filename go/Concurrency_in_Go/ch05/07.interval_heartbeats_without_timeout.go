package main

import (
	"testing"
	"time"
)

// 06과 동일
func DoWork(
	done <-chan interface{},
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWork(done, intSlice...)

	// 하트비트 신호가 들어올 때까지 대기
	<-heartbeat

	i := 0
	for r := range results {
		// DoWork에서 받은 값 r과
		// 전송을 시작했던 intSlice를 비교
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		// DoWork에서 작업 채널로 총 몇 개를 보냈는지 i로 저장
		i++
	}
}
