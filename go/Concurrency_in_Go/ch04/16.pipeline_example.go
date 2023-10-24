package main

import "fmt"

func main() {

	// 아래와 같은 방식과 함수를 이용한 방식의 차이점은
	// 입출력이 동시에 실행되는 컨텍스트에서 안전하다는 점.
	// 모든 stage는 동시에 입력 채널에 데이터가 들어오길 기다리며
	// 작업을 처리하여 출력 채널로 넘긴다

	// integers들을 채널 데이터 스트림으로 변경
	// stage들에게 integer를 주는 역할
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		// int 채널 생성
		intStream := make(chan int, len(integers))
		go func() {
			// 고루틴 종료시 int 채널 닫기
			defer close(intStream)
			// 인수로 들어온 모든 integer에 대해
			for _, i := range integers {
				select {
				// done 채널에 값이 들어오면 고루틴 종료
				case <-done:
					return
				// 모든 integer 차례대로 int 채널에 넣기
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		// int 데이터 스트림을 곱셈 처리한
		// int 데이터 스트림을 넘기는 채널 생성
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				// int 스트림 곱셈 처리 후 또 스트림
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				// // int 스트림 덧셈 처리 후 또 스트림
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
