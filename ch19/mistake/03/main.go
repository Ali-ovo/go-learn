package main

import (
	"fmt"
	"time"
)

func main() {
	var out []*int
	for i := 0; i < 3; i++ {
		out = append(out, &i)
	}

	for _, value := range out {
		fmt.Println(*value)
	}

	goodsID := []uint64{1, 2, 3, 4, 5}

	for _, value := range goodsID {
		go func() {
			fmt.Println("value:", value)
		}()
	}

	time.Sleep(3 * time.Second)
}
