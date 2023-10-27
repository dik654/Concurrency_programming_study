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
		// 채널들 생성
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		// 고루틴 실행
		go func() {
			// 고루틴 종료시 채널들 닫기
			defer close(heartbeat)
			defer close(results)

			// 하트비트 간격은 인수로 들어온 pulseInterval
			pulse := time.Tick(pulseInterval)
			// 작업의 간격은 pulseInterval * 2
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
					case <-done:
						return
					// 작업 중에도 pulseInterval에 도달할 경우
					case <-pulse:
						// 하트비트 보내기
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				// done 채널에 신호가 오면 고루틴 종료
				case <-done:
					return
				// 펄스 신호가 오면
				case <-pulse:
					// 하트비트 실행
					sendPulse()
				// pulseInterval * 2가 지났으면 하트비트
				// 또는 pulseInterval * 2가 지난 시간 채널로 보내기
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}
	done := make(chan interface{})
	// 10초 뒤에 done 채널 닫기
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
