package main

import (
	"fmt"
	"sync"
)

// 원자성이란 "동작하는 컨텍스트 내에서 나눠지거나 중단되지 않는 것"
// "동작하는 컨텍스트"이므로 다른 컨텍스트에서는 원자적이지 않을 수 있다(예. 사용자 앱의 컨텍스트와 운영체제의 컨텍스트)
// 또한 원자적인 연산을 조합한다고 원자성을 갖는 것은 아니다
// 한 예로 value++ 작업은
// 1. value를 읽어오기
// 2. value = value + 1
// 3. value를 쓰기
// 로 나뉘어 있다. 각각의 작업은 원자적이지만, 작업들이 합쳐진 value++은 원자적이지 않다.

func main() {
	// 뮤텍스를 이용하여 Lock을 먼저한 고루틴이 리소스를 독점하도록 하는 예
	// 이를 통해 메모리 접근을 동기화시켰다
	var memoryAccess sync.Mutex
	// 공유 리소스
	var value int
	go func() {
		memoryAccess.Lock()
		// 공유 리소스에 독점적으로 접근해야하는 부분(쓰기 critical section)
		value++
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	// 읽기 critical section
	if value == 0 {
		// 읽기 critical section
		fmt.Printf("the value is %v.\n", value)
	} else {
		// 읽기 critical section
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock()
}

// 위의 방법으로 데이터 접근에 대한 race는 해결했지만
// 어느 고루틴이 먼저 실행되는지 확실하지 않다는 race condition은 해결하지 못했다
// 또한 뮤텍스를 사용할 때, 또 다른 문제들이 생길 수 있다
// - deadlock
// - livelock
// - starvation
