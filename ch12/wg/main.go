package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()

	// wait group 用于 goroutine 执行等待 Add 配合 Done配套使用
}
