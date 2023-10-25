package main

import "fmt"

func main() {
	bridge := func(
		done <-chan interface{},
		// 읽기 전용 채널을 전달하는 읽기전용 채널
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})
		// 고루틴 생성
		go func() {
			defer close(valStream)
			for {
				// 전달 받은 채널을 연결할 bridge 채널 생성
				var stream <-chan interface{}
				select {
				// 읽기 전용 채널을 전달하는 채널에서 신호가 왔고
				case maybeStream, ok := <-chanStream:
					// 채널에 더 이상 읽을 값이 없고
					// 채널이 닫혀있다면
					if ok == false {
						// 고루틴 종료
						return
					}
					// bridge 채널에 전달받은 채널 연결
					stream = maybeStream
				case <-done:
					return
				}
				// orDone 패턴을 통해 bridge 채널을 통해 데이터 읽기
				for val := range orDone(done, stream) {
					select {
					// 읽은 데이터는 출력 채널인 valStream 채널로 전달
					case valStream <- val:
					// done 신호가 올 경우 orDone 고루틴에서 done 처리
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	genVals := func() <-chan <-chan interface{} {
		// 읽기전용 채널을 전달하는 채널 생성
		chanStream := make(chan (<-chan interface{}))
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 읽기전용 채널을 전달하는 채널 닫기
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				// 임시 채널 하나를 생성하여
				stream := make(chan interface{}, 1)
				// for문의 i에 해당하는 숫자를 임시 채널에 넣고
				stream <- i
				// 임시 채널 닫기
				close(stream)
				// 닫힌 채널 전달(i값은 들어있어서 읽을 수는 있음)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	// 채널을 스트림으로 전달받고
	// bridge 내부 임시 채널을 이용해
	// 정상적인 채널이 들어왔으면 값을 출력 채널로 전달하고
	// 닫힌 채널이 들어왔으면 고루틴 종료

	// 숫자가 담긴 닫힌 채널을 생성해 전달하는 채널을 리턴하는 genVals()
	// 채널을 bridge로 넘겨 임시 채널을 이용해 닫힌 채널의 값을 읽는 과정
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
