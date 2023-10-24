package main

import "fmt"

func main() {
	// doWork함수는 strings라는 이름의 string 읽기 채널을 받아서,
	// any 읽기 채널을 반환한다
	doWork := func(strings <-chan string) <-chan interface{} {
		// completed라는 any 채널을 만들고
		completed := make(chan interface{})
		// 고루틴을 실행시킨다
		go func() {
			// 고루틴 종료시
			defer fmt.Println("doWork exited.")
			// any 채널을 닫는다
			defer close(completed)
			// strings 채널이 닫힐 때까지 계속
			// strings 채널에 들어온 string을 읽어서
			for s := range strings {
				// 콘솔에 찍는다
				fmt.Println(s)
			}
		}()
		// 고루틴을 생성시키면서 completed 채널을 리턴한다
		return completed
	}
	// strings 채널로 nil 채널을 보내
	// 아무 데이터도 strings 채널로 가지않는다
	// 즉, 고루틴이 아무런 동작은 하지않고 계속 메모리를 차지한다
	doWork(nil)
	fmt.Println("Done.")
}
