package main

import "fmt"

// 아래 코드의 의도는 생성한 고루틴이 data 변수를 변경시키는 동시에
// main 고루틴이 data 변수를 읽는 과정에서
// 순서상의 문제로 상황에 따라 결과가 다르게 나올 수 있다는 점이다

// 하지만 쉘 스크립트 while true; do go run main.go; done으로 여러번 실행시켜본 결과
// 0 이외의 결과는 확인하지 못했다
// 왜냐하면 일반적으로 main 고루틴이 먼저 실행될 확률이 높아서이다.
func main() {
	// 초기값 0
	var data int
	// 새로 생성한 고루틴은 CPU를 얻어 실제로 실행되기까지 시간이 걸릴 수 있기 때문
	go func() {
		data++
	}()

	// 초기값이 0이고, 대부분의 상황에서 main 고루틴이 새로 생성한 고루틴보다 빨리 동작하므로
	// if data == 0를 통과하고 변수 값인 0을 리턴할 확률이 매우 높다
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	} else {
		fmt.Printf("nothing")
	}

	// 하지만 항상 같은 결과를 갖는 것이 보장되지 않으므로 주의해야한다
}
