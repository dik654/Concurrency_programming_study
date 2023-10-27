package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

// 에러 타입에 맞게 래핑
func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

// 에러 메세지 getter
func (err MyError) Error() string {
	return err.Message
}

// 상속 패턴
// LowLevelErr 타입이 error의 메서드 상속
// 에러가 어느 패키지에서 일어났는지 stack trace를 위해 선언
type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	// path의 지정된 파일의 정보 info 변수로 가져오기
	info, err := os.Stat(path)
	// 가져오는 도중 실패했다면
	if err != nil {
		// 에러를 MyError타입에 맞춰 추적이 가능하도록 하고
		// import한 모듈에서 생긴 에러이니 해당 모듈 에러 타입으로 선언하여 리턴
		return false, LowLevelErr{(wrapError(err, err.Error()))}
	}
	// 유닉스 파일모드가 찾고있는 어떤 한 부분이 실행 권한을 갖고있는지 체크
	return info.Mode().Perm()&0100 == 0100, nil
}

// IntermediateErr 타입이 error의 메서드 상속
// 에러가 어느 패키지에서 일어났는지 stack trace를 위해 선언
type IntermediateErr struct {
	error
}

func runJob(id string) error {
	// 명령어 경로
	const jobBinPath = "/bad/job/binary"
	// 실행이 가능하다면
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return err
	} else if isExecutable == false {
		// 실행이 불가하다면 wrapError 구조에 맞춰서 에러 생성
		// 그러나 import한 모듈이 아닌, 내 작업내에서 일어난 에러이므로
		// 따로 에러 타입을 선언하지는 않는다
		return wrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) // <3>
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	// /bad/job/binary 명령어를 --id=1 인수와 함께 실행
	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		// runJob에서 타입 래핑을 넘어가서 타입 체크에 실패
		// [1] There was an unexpected issue; please report this as a bug를 리턴하게 됨
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		// 가져온 에러 콘솔에 뿌리기
		handleError(1, err, msg)
	}
}
