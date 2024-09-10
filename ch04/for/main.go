package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 3; i++ {
		println(i)
	}

	var i int
	for i < 3 {
		time.Sleep(1 * time.Second)
		fmt.Println(i)
		if i == 1 {
			break
		}
		i++
	}

	var sum int
	for i := 1; i < 100; i++ {
		sum += i
	}
	fmt.Println(sum)

	// for range 对于数组 字符串 切片等处理
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}

	name := "ali2333还有中文"
	for i, val := range name {
		fmt.Printf("%c ", name[i])
		fmt.Printf("index: %d, value: %c \n", i, val)
	}
}
