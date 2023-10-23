package main

import "fmt"

func main() {
	chanOwner := func() <-chan int {
		// chanOwner 어휘 범위 내에서 채널을 인스턴스화
		// 버퍼 크기 5의 int 채널 생성
		results := make(chan int, 5)
		// 고루틴 생성
		go func() {
			// results 채널 닫기
			defer close(results)
			for i := 0; i <= 5; i++ {
				// 채널에 0 ~ 5 넣기
				// 6개를 넣어서 버퍼 크기가 넘어서면 쓰기 과정이 블로킹
				results <- i
			}
		}()
		return results
	}

	// 채널 읽기(<-chan)만 받아서 읽기만 할 수 있도록 제한
	// 채널을 인수로 받는 컨슈머 함수 선언
	consumer := func(results <-chan int) {
		// 채널에서 데이터를 꺼내서
		for result := range results {
			// 차례대로 화면에 뿌리기
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}

	// 채널을 복사하여 results에 저장
	results := chanOwner()
	// 복사한 채널 소비하기
	consumer(results)
}
