package main

import (
	"fmt"
	"math/rand"
)

func main() {
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
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				// 인수로 들어온 함수를 실행시킨 결과값
				// 계속 채널로 넘기기
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	// 랜덤 숫자를 리턴하는 함수 생성
	rand := func() interface{} { return rand.Int() }

	// range는 채널에서 데이터를 읽어와 변수에 저장시킨다
	// 이는 채널이 닫힐 때까지 반복된다
	// 생성한 함수를 반복적으로 넘기고 10번 다음 채널로 넘기기
	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}
