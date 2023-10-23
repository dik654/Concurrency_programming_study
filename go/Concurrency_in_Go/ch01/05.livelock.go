package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// livelock
// 프로그램이 열심히 동작하고 있지만 프로그램 상태에 진척이 없는 경우
func main() {
	// 특정 조건이 만족할 때까지 고루틴들을 대기하도록 만들 수 있음
	cadence := sync.NewCond(&sync.Mutex{})
	// 0.1초 마다 고루틴들 깨우기
	go func() {
		// 1밀리초마다 뮤텍스 타이밍 맞추기
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	// 서로 동시에 움직이도록 락을 걸고 고루틴들 대기시키기
	takeStep := func() {
		// 고루틴들의 공유 리소스 동시 접근 제한
		cadence.L.Lock()
		// 고루틴들 대기 상태로 변경, Broadcast()를 받으면 다시
		cadence.Wait()
		// 동시 접근 제한 해제
		cadence.L.Unlock()
	}

	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
		// 콘솔에 이동하려는 방향 리턴
		fmt.Fprintf(out, " %v", dirName)
		// 반드시 해당 방향으로 +1 작업을 하고
		atomic.AddInt32(dir, 1)
		// 고루틴들 대기
		takeStep()
		// 한쪽이 이동을 먼저 성공해서 dir값이 1이 됐다면
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			// 이동 성공, true 리턴
			return true
		}
		// tryDir()가 이미 한 번 성공했거나
		// 늦게 이동한 나머지 한 명은 dir이 2가 되어
		takeStep()
		// dir를 1로 만듦
		atomic.AddInt32(dir, -1)
		// 그리고 false 리턴
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		// 이동하려고한 방향 리턴
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			// 왼쪽 먼저 이동 후, 오른쪽 이동 시도
			if tryLeft(&out) || tryRight(&out) {
				// 한번이라도 성공하면 이동 시도 탈출
				return
			}
		}
		// 이동에 실패하면 콘솔에 뿌리기
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
