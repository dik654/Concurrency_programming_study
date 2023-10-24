package main

import (
	"fmt"
	"time"
)

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	// or 함수는 여러 채널을 인수로 받는다(...)
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		// 인수 개수에 따라 다른 동작
		switch len(channels) {
		// 인수가 없을 경우
		case 0:
			return nil
		// 1개일 경우
		case 1:
			// 인수로 받은 그 채널 리턴
			return channels[0]
		}

		// 종료 신호 채널 생성
		orDone := make(chan interface{})
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 orDone 채널 닫기
			defer close(orDone)

			// 인수로 받은 채널 개수가
			switch len(channels) {
			// 2개일 경우
			case 2:
				// 0, 1번 채널에서 데이터 들어올 때까지 대기
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// 2개보다 많다면
			default:
				select {
				// 0, 1, 2번 채널에서 데이터가 들어올 때까지 대기
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				// 채널을 3개씩 나누기
				// 나머지는 재귀적으로 다른 3개 이하 채널에 데이터가 들어오길 기다리는 고루틴을 생성하도록함
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		// 이렇게 신호를 대기하는 재귀적으로 생성된 고루틴들을
		// 한번에 종료할 수 있도록 orDone 채널 리턴
		return orDone
	}

	// 시간을 인수로 받는다
	sig := func(after time.Duration) <-chan interface{} {
		// any 채널 생성
		c := make(chan interface{})
		// 고루틴 생성
		go func() {
			// 인수로 받은 시간만큼 기다린 후
			// c 채널 닫기
			defer close(c)
			time.Sleep(after)
		}()
		// 생성한 c 채널 리턴
		return c
	}
	// 시작 시간 저장
	start := time.Now()
	// orDone 채널에서 종료 신호를 받을 때까지 대기
	<-or(
		// 인수로 받은 시간 후 종료 신호를 뿌리는 c채널들을 or함수의 인수로 넣는다
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	// 가장 적은 시간이 1초이므로 1초 뒤에 c 채널에 채널닫힘 신호를 받고
	// or 함수에서 생성된 모든 고루틴들이 종료된다
	// done after 1.000813447s
	fmt.Printf("done after %v", time.Since(start))
}
