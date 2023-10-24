package main

import "fmt"

func main() {

	repeat := func(
		done <-chan interface{},
		// 모든 타입의 입력값 받기
		values ...interface{},
	) <-chan interface{} {
		// 채널 생성
		valueStream := make(chan interface{})
		// 고루틴 생성
		go func() {
			// 고루틴 종료 채널 닫기
			defer close(valueStream)
			// 입력에 대해 한번만 뿌리는게 아닌 반복적으로 뿌리기
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					// 채널에 넣기
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				// valueStream에서 값을 읽어와서(<-valueStream)
				// takeStream 채널로 넘기기 (<-)
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	// 1을 반복적으로 넣고 10번 채널에서 가져와서 콘솔에 뿌리기
	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}
