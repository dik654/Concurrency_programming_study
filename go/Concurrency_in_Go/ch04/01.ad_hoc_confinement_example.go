package main

import "fmt"

// 애드혹 제한이란 일반적인 규칙이나 원칙이 아니라
// 특정 문제나 상황을 해결하기 위해 임시로 설정된 규칙
func main() {
	// int 타입 4개를 담을 수 있는 슬라이스 생성
	data := make([]int, 4)

	// 함수의 인수로 채널 넘기기
	// 슬라이스를 넘기는게 아니라 int를 넘기는 규칙으로 짜는 예
	// 실수 등에 의해 이런 제한이 깨질 수 있다
	loopData := func(handleData chan<- int) {
		// 함수 종료시 채널 닫기
		defer close(handleData)
		// 슬라이스 내 데이터 모두 채널로 넘기기
		for i := range data {
			// handleData에서 슬라이스를 쓸 수도있지만
			handleData <- data[i]
		}
	}

	// int 채널 생성
	handleData := make(chan int)
	// 생성한 채널을 인수로 고루틴 생성
	go loopData(handleData)

	// 채널에 있는 데이터 차례대로 콘솔에 뿌리기
	for num := range handleData {
		fmt.Println(num)
	}
}
