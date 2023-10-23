package main

import (
	"fmt"
	"sync"
	"time"
)

// CoffmanDeadlocks(deadlock이 발생하는 조건)
// 아래 조건이 모두 참일 때 데드락
// Mutual Exclusion - 동시에 실행되는 프로세스들이, 같은 시점에 리소스에 대한 독점을 요구하는 경우
// Wait For Condition - 이미 할당된 자원을 갖고 있으면서, 다른 자원이 풀리길 기다리는 경우
// No Preemption - 이미 할당된 자원을 강제로 뺏을 수 없는 경우
// Circular Wait - A가 *a를 보유하고 *b의 선점이 풀리길 대기, B가 *b를 보유하고 *a의 선점이 풀리길 대기하는 경우

func main() {
	type value struct {
		mu    sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		// 함수 종료시 wg 완료
		defer wg.Done()
		// 첫 번째 인수 락 걸고
		v1.mu.Lock()
		defer v1.mu.Unlock()

		// 2초 대기 후
		time.Sleep(2 * time.Second)
		// 두 번째 인수 락 건 뒤
		v2.mu.Lock()
		defer v2.mu.Unlock()

		// 두 변수에 접근
		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	// 고루틴 두개 생성하고 대기하기 위해 wg 2개 추가
	wg.Add(2)
	// a에 대해서 락 건 뒤
	// 2초 뒤 b에 대해서 락 걸기 시도
	go printSum(&a, &b)
	// b에 대해서 락 건 뒤
	// 2초 뒤 a에 대해서 락 걸기 시도
	go printSum(&b, &a)
	// 두 고루틴에서 a, b가 이미 락이 걸려있어 2초 뒤에 락을 걸려고 해도 걸 수 없다
	// 따라서 printSum 고루틴 2개가 종료되어야 동작하는 아래 코드에는 도달할 수 없다
	// fatal error: all goroutines are asleep - deadlock!
	wg.Wait()
}
