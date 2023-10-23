package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		// 현재시간부터 1초될 때까지
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			// 3 나노초 독점
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			// 1초씩 3번 독점
			// starvation 상태에 빠질 확률이 높음
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}
	// wait group 생성
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	// 두 고루틴에서 모두 wg.Done이 일어나도록 대기
	wg.Wait()
}
