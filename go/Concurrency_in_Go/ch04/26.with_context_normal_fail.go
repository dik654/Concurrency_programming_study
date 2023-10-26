package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// done 채널 대신 컨텍스트를 받는다
func locale(ctx context.Context) (string, error) {
	select {
	// 컨텍스트 취소시
	case <-ctx.Done():
		// context.Canceled 리턴
		return "", ctx.Err()
	// 1분 뒤
	case <-time.After(1 * time.Minute):
	}
	// EN/US 리턴
	return "EN/US", nil
}

// 컨텍스트를 받아서
func genGreeting(ctx context.Context) (string, error) {
	// 인수로 받은 부모 ctx가 cancel되면 같이 취소되는
	// 1초 뒤 취소되는 자식 컨텍스트 생성
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	// 1초 뒤에 취소되는 컨텍스트를 넘겼는데,
	// 1분 뒤에 locale 정보를 받을 수 있으므로
	// context deadline exceeded
	switch locale, err := locale(ctx); {
	// 컨텍스트에 에러가 담겨있을 경우 컨텍스트 에러 그대로 넘기기
	case err != nil:
		return "", err
	// EN/US일 경우 hello 넘기기
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	// 위에서 내려온 컨텍스트 에러 넘기기
	case err != nil:
		return "", err
	// 셩공시 goodbye 리턴
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func printGreeting(ctx context.Context) error {
	// hello 받아서
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	// 콘솔에 뿌리기
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	// goodbye 받아서
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	// 콘솔에 뿌리기
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func main() {
	var wg sync.WaitGroup
	// done역할을 포함하는 cancel 컨텍스트 생성
	ctx, cancel := context.WithCancel(context.Background())
	// 메인 함수 종료시 컨텍스트 취소
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// 컨텍스트를 넘겨서 hello받기
		if err := printGreeting(ctx); err != nil {
			// 에러가 났을 경우 에러 상황 콘솔에 프린트
			fmt.Printf("cannot print greeting: %v\n", err)
			// 후 컨텍스트 취소
			// 1초 뒤 타임아웃되는 컨텍스트 때문에 에러가 발생하여
			// 메인 함수에서 선언한 부모 컨텍스트가 취소된다
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			// greeting 타임아웃으로 부모 컨텍스트가 취소되어
			// context canceled
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()
	// wg이 0 될 때까지 대기
	wg.Wait()
}
