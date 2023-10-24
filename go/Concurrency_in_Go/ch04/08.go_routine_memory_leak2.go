package main

import (
	"fmt"
	"math/rand"
)

func main() {
	newRandStream := func() <-chan int {
		// 채널 생성
		randStream := make(chan int)
		// 고루틴 생성
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			// 고루틴 종료시 randStream 채널 닫기
			defer close(randStream)
			for {
				// 난수를 생성해서 randStream 채널에 넘기기
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	// 채널을 닫지않고 채널에서 값 3개만 꺼내서 콘솔에 뿌리고 종료
	// memory leak이 일어난다
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}
