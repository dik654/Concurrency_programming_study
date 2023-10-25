package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// fan out - 파이프라인 입력처리를 위해 여러 고루틴을 실행
// fan in - 여러 결과를 하나의 채널로 가져오기

// A라는 작업이 있을 때 fan out을 적용하기 적합한 상황은
// - A 작업은 시간이 오래걸리지만,
// - A 작업 이전의 작업의 결과나 계산에 의존하지 않는 경우 작힙히디

func main() {
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))

		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	rand := func() interface{} { return rand.Intn(50000000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	// 랜덤 정수 가져오기
	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")

	// 내 컴퓨터의 cpu 개수
	numFinders := runtime.NumCPU()
	//
	finders := make([]<-chan int, numFinders)
	// cpu 개수만큼 소수 찾는 고루틴 생성
	for i := 0; i < numFinders; i++ {
		// 고루틴 채널 finders에 저장
		finders[i] = primeFinder(done, randIntStream)
	}

	// 여러 고루틴이 소수를 찾는다
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	// 10개를 가져오는데 걸린 총 시간
	fmt.Printf("Search took: %v", time.Since(start))
}
