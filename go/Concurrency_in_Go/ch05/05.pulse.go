package main

import (
	"fmt"
	"math/rand"
)

func main() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		// 하트비트용 채널 생성
		heartbeatStream := make(chan interface{}, 1) // <1>
		// 작업 데이터를 보낼 채널 생성
		workStream := make(chan int)
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 채널 모두 닫기
			defer close(heartbeatStream)
			defer close(workStream)

			// 10번 반복
			for i := 0; i < 10; i++ {
				// 랜덤 값 작업 전에 하트비트 날리기
				select {
				// 하트비트 채널이 가득찼거나 닫힌 상태가 아니라면 struct{}{}를 전송
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				// done 채널에 신호가 들어왔다면 고루틴 종료
				case <-done:
					return
				// 0 ~ 9 사이의 랜덤 값 작업 채널로 넘기기
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	// done 채널 생성
	done := make(chan interface{})
	// main 고루틴 종료시 done 채널 닫기
	defer close(done)

	// 10번 랜덤 값 뿌리는 고루틴 실행
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
