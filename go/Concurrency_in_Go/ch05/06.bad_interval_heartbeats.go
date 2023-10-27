package main

import (
	"testing"
	"time"
)

func DoWork(
	done <-chan interface{},
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)
		// 2초간 중지하여 작업이 있는 상황 mocking
		time.Sleep(2 * time.Second)

		// 인수로 들어온 숫자 슬라이스들에 대하여
		for _, n := range nums {
			select {
			// 하트비트 채널에 버퍼가 남아있거나 닫히지 않았다면 struct{}{}를 보내고
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			// done 채널에 값이 들어온 경우 고루틴 종료
			case <-done:
				return
			// 인수로 들어온 숫자 intStream 채널로 전송
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
	// 숫자 슬라이스를 채널로 보내는 DoWork 고루틴 실행
	_, results := DoWork(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			// 넣은 순서대로 DoWork에서 나오지 않는다면
			if r != expected {
				// 테스트 에러 콘솔에 에러 뿌리기
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		// DoWork의 작업이 1초가 넘어간다면 테스트 실패
		case <-time.After(1 * time.Second):
			t.Fatal("test timed out")
		}
	}
}
