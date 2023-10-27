package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			// 2번의 신호를 보내고 채널이 닫히지 않은 비정상적인 상황 연출

			// 하트비트는 인수로 들어온 pulseInterval만큼
			pulse := time.Tick(pulseInterval)
			// 작업은 그 2배의 간격으로 설정
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				// pulseInterval이 지났으면
				case <-pulse:
					// 하트비트 채널로 신호 보내기
					sendPulse()
				// 고루틴 시작으로부터 pulseInterval의 2배가 지난 시간이 지났다면
				case r := <-workGen:
					// results로 workGen에
					// 고루틴 시작으로부터 pulseInterval의 2배가 지난 시간 저장
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	// done 채널 생성
	done := make(chan interface{})
	// 10초 뒤 인수로 보낸 함수 실행시키기
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		// 하트비트 채널에 신호가 들어온 경우
		case _, ok := <-heartbeat:
			// 채널에 데이터가 없고 닫혔다면
			if ok == false {
				// main 함수 종료
				return
			}
			fmt.Println("pulse")
		// 작업 채널로 고루틴 시작으로부터 pulseInterval의 2배가 지난 시간이 넘어왔다면
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r)
		// 하트비트가 없고, results에도 시간 값이 안넘어오고 2초가 지났다면(즉, 비정상적인 상황이라면)
		case <-time.After(timeout):
			// 고루틴에 문제가 있음을 콘솔로 뿌리고 main함수 종료
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}
