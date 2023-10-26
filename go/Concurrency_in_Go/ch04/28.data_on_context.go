package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "abc123")
}

func ProcessRequest(userID, authToken string) {
	// 컨텍스트에
	// key: userID,
	// value: userID 저장
	ctx := context.WithValue(context.Background(), "userID", userID)
	// 컨텍스트에
	// key: authToken,
	// value: authToken 저장
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (%v)",
		// 컨텍스트 내 key인 userID에 저장된 값 리턴
		ctx.Value("userID"),
		// 컨텍스트 내 key인 authToken에 저장된 값 리턴
		ctx.Value("authToken"),
	)
}
