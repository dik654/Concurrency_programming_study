package main

import "fmt"

func main() {
	type foo int
	type bar int

	m := make(map[interface{}]int)
	m[foo(1)] = 1
	m[bar(1)] = 2

	// map[1:2 1:1]
	// 같은 int 타입이라도 선언한 타입명이 다르다면
	// 따로 저장된다
	fmt.Printf("%v", m)
}
