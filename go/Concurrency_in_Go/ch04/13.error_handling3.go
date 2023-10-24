package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 에러 이유와 응답을 담은 Result 타입
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		// results 채널 생성
		results := make(chan Result)
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 results 채널 닫기
			defer close(results)

			// 인수로 들어온 모든 urls들에 대해서
			for _, url := range urls {
				var result Result
				// GET 요청 날리기
				resp, err := http.Get(url)
				// 해당 응답과 에러를 에러 객체로 만들어 result에 담기
				result = Result{Error: err, Response: resp}
				// 아래의 조건에 해당할 때까지 대기
				select {
				// done채널에 신호가 들어오면 고루틴 종료
				case <-done:
					return
				// results 채널에 Result 객체 보내기
				case results <- result:
				}
			}
		}()
		return results
	}

	// done 채널 만들기
	done := make(chan interface{})
	// done 채널 닫기
	defer close(done)

	errCount := 0
	// 테스트 케이스
	urls := []string{"a", "https://www.google.com", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		// 만약 GET요청시 에러가 있다면
		if result.Error != nil {
			// 에러 내용 뿌리고
			fmt.Printf("error: %v", result.Error)
			// 발견한 에러 숫자 증가시키고
			errCount++
			// 만약 에러가 3개 이상이면 for문 종료
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			// 아니라면 다음 케이스로 넘어가기
			continue
		}
		// 에러가 없다면 응답 내용 뿌리기
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
