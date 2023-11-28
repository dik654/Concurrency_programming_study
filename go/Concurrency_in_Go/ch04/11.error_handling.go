package main

import (
	"fmt"
	"net/http"
)

func main() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		// http 응답이 이동하는 채널 생성
		responses := make(chan *http.Response)
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 응답 이동 채널 닫기
			defer close(responses)
			// 인수로 들어온 모든 url들에 대해
			for _, url := range urls {
				// GET 요청을 날리고
				resp, err := http.Get(url)
				// 요청시 에러가 나타나면
				if err != nil {
					// 콘솔에 에러사항을 뿌리고
					fmt.Println(err)
					// 다음 url로 넘어가기
					continue
				}
				// 아래 조건이 만족할 때까지 대기
				select {
				// done 채널에 신호가 들어오거나
				case <-done:
					// 그러면 고루틴 종료
					return
				// 받은 응답이 있으면 responses 채널에 넘기기
				case responses <- resp:
				}
			}
		}()
		return responses
	}
	// done 채널 만들기
	done := make(chan interface{})
	// 함수 종료시 고루틴 종료
	defer close(done)

	// 보낼 url들
	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}
