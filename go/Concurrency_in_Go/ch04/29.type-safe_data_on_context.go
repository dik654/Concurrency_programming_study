package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	// 자동으로 ctxKey타입이 됨
	ctxAuthToken
)

// 컨텍스트에서 값을 가져올 때 타입을 고려하지 않아도 되도록 getter함수 선언
func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

// 컨텍스트에서 값을 가져올 때 타입을 고려하지 않아도 되도록 getter함수 선언
func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

// 컨텍스트에 값 넣는 과정
func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

// 리액트 상태관리와 유사하게
// 만약 HandleResponse가 이 패키지가 아닌 다른 B라는 패키지안에 있다면
// ProcessRequest가 HandleResponse를 호출하려면 B패키지를 호출해야하고
// HandleResponse가 UserID(), AuthToken()을 호출하려면 이 패키지를 호출해야하므로
// 순환 의존성 문제가 생긴다.

// 따라서 한 방향으로만 의존성을 갖도록 구조를 짜야한다
func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}
