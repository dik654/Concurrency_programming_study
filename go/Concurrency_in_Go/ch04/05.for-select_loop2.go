package main

import "os"

func main() {
	done := make(chan os.Signal, 1)

	// 인터럽트가 올 때까지 무한정 대기하는 케이스
	for {
		select {
		case <-done:
			return
		default:
		}
	}
}
