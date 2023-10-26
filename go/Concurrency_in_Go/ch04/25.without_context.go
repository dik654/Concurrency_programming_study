package main

import (
	"fmt"
	"sync"
	"time"
)

// context는 done 채널에 있어
// 취소 이유, deadline등의 추가 정보를 전달하기 위해 추가되었다

func locale(done <-chan interface{}) (string, error) {
	// 아래 조건에 해당할 때까지 대기
	select {
	// done에 신호가 들어오면
	case <-done:
		// canceled 에러를 리턴하고 함수 종료
		return "", fmt.Errorf("canceled")
	// time.After메서드는 1분 후 새로운 채널을 생성하여 리턴
	// 그 채널에는 1분 뒤의 현재 시간를 전송

	// 해당 코드는 채널에 꺼내올 값이 있을 때 실행되므로
	// "시간을 전달하는 채널"이 전달 되었을 때
	// <-time.After(1 * time.Minute)의 채널에 "시간을 전달하는 채널"이 전달되었을 때
	// 아래의 case가 만족되어 실행
	case <-time.After(1 * time.Minute):
	}
	// 종료시 "EN/US" 리턴
	return "EN/US", nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	// 1분 뒤 EN/US를 받는다
	switch locale, err := locale(done); {
	// done에 의해 취소되었을 경우 위에서 내려온 에러 전달
	case err != nil:
		return "", err
	// hello를 리턴한다
	case locale == "EN/US":
		return "hello", nil
	}
	// EN/US가 아니라면 에러 발생
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	// 1분 뒤 EN/US를 받는다
	switch locale, err := locale(done); {
	// done에 의해 취소되었을 경우 위에서 내려온 에러 전달
	case err != nil:
		return "", err
	// goodbye를 리턴한다
	case locale == "EN/US":
		return "goodbye", nil
	}
	// EN/US가 아니라면 에러 발생
	return "", fmt.Errorf("unsupported locale")
}

func printGreeting(done <-chan interface{}) error {
	// 1분이 지나고 locale에서 EN/US를 받고
	// genGreeting에서 hello를 받아서 greeting 변수에 저장
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	// hello world! 콘솔에 뿌리기
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	// 1분이 지나고 locale에서 EN/US를 받고
	// genGreeting에서 hello를 받아서 greeting 변수에 저장
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	// goodbye world! 콘솔에 뿌리기
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	// wait group을 추가하고
	wg.Add(1)
	// 고루틴 실행
	go func() {
		// 종료시 wg--
		defer wg.Done()
		//
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	// 두 고루틴은 race condition 때문에
	// 어떤 함수가 먼저 결과를 리턴할 지 알 수 없다
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	// wg이 0이 될 때까지 대기
	wg.Wait()
}
