package main

func main() {

	// 다른 채널을 사용할 때,
	// 그 채널이 취소된 상태일 때의 동작을 확실하는 패턴
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 valStream 채널 닫기
			defer close(valStream)
			for {
				select {
				// done 채널에 신호가 들어오면 고루틴 종료
				case <-done:
					return
				// 인수로 들어온 채널에 값이 들어오면
				case v, ok := <-c:
					// 채널이 취소된 상태일 경우를 확실히 하는 부분
					// 채널이 닫혀있는 상태라면
					if ok == false {
						// 고루틴 종료
						return
					}
					// 케이스 중 하나가 나올 때까지 대기
					select {
					// 들어온 값 valStream에 보내기
					case valStream <- v:
					// done에 들어온 값 읽기
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	// 장황한 select case문을 캡슐화하여 사용
	for val := range orDone(done, myChan) {
		// val로 무언가
	}
}
