package main

import "fmt"

func main() {
	ints := []int{1, 2, 3, 4}
	// 스트림 관점에서 작업을 생각하기 위해 단일 int에 대한 계산으로 stage를 변경
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}

	add := func(value, additive int) int {
		return value + additive
	}

	// #####################################################################

	// 절차적으로 작성됐지만 파이프라인의 고수준 결합을 사용할 수 있다
	for _, v := range ints {
		fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	}
}
