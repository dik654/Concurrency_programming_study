package main

import (
	"fmt"
	"math/rand"
	"time"
)

// fan out - 파이프라인 입력처리를 위해 여러 고루틴을 실행
// fan in - 여러 결과를 하나의 채널로 가져오기

// A라는 작업이 있을 때 fan out을 적용하기 적합한 상황은
// - A 작업은 시간이 오래걸리지만,
// - A 작업 이전의 작업의 결과나 계산에 의존하지 않는 경우 작힙히디

func main() {
	rand := func() interface{} { return rand.Intn(5000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	// primeFinder 단계에서 한정없이 소수를 찾아서 채널로 넘기면
	// take에서 10개까지 소수를 가져온다
	// 하나의 primeFinder에서 소수를 찾으므로 속도가 느리다
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	// 10개를 가져오는데 걸린 총 시간
	fmt.Printf("Search took: %v", time.Since(start))
}
