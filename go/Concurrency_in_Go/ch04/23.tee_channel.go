package main

import "fmt"

func main() {
	// 하나의 채널에서 들어온 데이터를 복사하여 여러 채널로 나눠주는 패턴
	//         o		in
	//         |
	//         |
	//         v
	//         o		orDone
	//        / \
	//       |   |
	//       v   v
	//       o   o		out1, 2

	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (_, _ <-chan interface{}) {
		// 채널 2개 생성
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		// 고루틴 생성
		go func() {
			// 고루틴 종료시 채널 닫기
			defer close(out1)
			defer close(out2)
			// in 채널이 닫혔을 때의 상황을 확실히하면서
			// valStream 채널로 데이터 받기
			for val := range orDone(done, in) {
				// 아래에서 out = nil 때문에 채널이 무시되기 때문에
				// orDone의 채널에 데이터가 들어올 때마다 out 채널들을 초기화 시켜준다
				var out1, out2 = out1, out2
				// 두 채널 모두에 전송하기 위해 2번 반복
				for i := 0; i < 2; i++ {
					// 처음에는 아래 2개의 케이스를 모두 만족하므로 랜덤으로 선택됨
					select {
					case <-done:
					// 채널에 데이터를 썼다면
					case out1 <- val:
						// 데이터를 받은 채널을 nil로 설정하여
						// 데이터가 들어와도 무시하도록 설정
						out1 = nil
						// 마찬가지로 동작하여 양쪽 채널에 같은 데이터가 들어가게 된다
						// out 채널들은 nil로 설정되어 이 함수는 첫 1회만 사용 가능하다
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
