package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	// printData 밖에서 선언된 data 슬라이스를
	// 다른 클로저인 printData내에서 포인터로 가져와서 직접 사용하는 것이 아닌
	// 복사해서 사용하여 제한
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}
