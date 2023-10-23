package main

import (
	"fmt"
	"time"
)

// 아래 코드를 여러 번 실행해본 결과
// 아래 코드는 기존 01과 달리 main 고루틴을 1초 대기시켜
// 고루틴의 로직이 실행될 시간을 주어 if data == 0절이 통과되지 않았다.

// 그러나 실제 환경에서는 대기하는 시간이 밀리, 마이크로 초 단위일 확률이 높다.
// 즉, 대기하는 시간이 길 수록 다른 가능성들이 낮아진 것일 뿐
// 로직이 확정적이 된 것은 아니다.

func main() {
	var data int
	go func() {
		data++
	}()
	time.Sleep(1 * time.Second)
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	}
}
