package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	done := make(chan os.Signal, 1)
	stringStream := make(chan string, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			if s, ok := <-stringStream; ok {
				println(s)
			} else {
				break
			}
		}
	}()

	// 순회할 수 있는 걸 채널의 값으로 변환
	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
}
