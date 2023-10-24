package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(
		// 종료 신호 채널
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		// any 타입의 종료 신호 채널 생성
		terminated := make(chan interface{})
		// 고루틴 실행
		go func() {
			defer fmt.Println("doWork exited.")
			// terminated 채널 닫기
			defer close(terminated)
			for {
				// 아래 케이스에 해당할 때까지 대기
				select {
				// strings 채널에 값이 들어오면 s에 저장한 뒤
				case s := <-strings:
					// 콘솔에 뿌리기
					fmt.Println(s)
				// done 채널에 신호가 들어오면
				case <-done:
					// 고루틴 종료
					return
				}
			}
		}()
		return terminated
	}

	// any 타입의 종료 신호 채널 생성
	done := make(chan interface{})
	// 종료될 때 terminated 채널을 닫는 고루틴 생성
	terminated := doWork(done, nil)

	// 고루틴 생성
	go func() {
		// 1초 뒤
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		// done 채널을 닫는다
		// done 채널을 닫을 때 doWork에서 생성한 고루틴이 종료되므로
		// terminated 채널도 닫힌다
		close(done)
	}()

	// 채널에 값이 들어올 때까지 대기
	// 채널이 닫히면 채널 타입에 해당하는 초기값이 반환
	// interface 타입이므로 nil이 terminated 채널로 넘어감
	<-terminated
	fmt.Println("Done.")
}
