package main

import "fmt"

func main() {
	ints := []int{1, 2, 3, 4}
	// Stage들
	// 일괄적으로 곱셈, 덧셈을 처리
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	// #####################################################################

	// 각 stage를 수정하지않고 단계들을 쉽게 합칠 수 있다
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}

	// 아래와 같이 스트림을 절차적으로 작성한 경우 고수준에서 고려하기 힘들다
	for _, v := range ints {
		fmt.Println(2 * (v*2 + 1))
	}

}
