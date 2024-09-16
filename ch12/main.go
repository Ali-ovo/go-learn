package main

import (
	"fmt"
	"time"
)

func asyncPrint() {
	fmt.Println("ali")
}

func main() {

	go asyncPrint()

	go func() {
		fmt.Println("hello")
	}()

	fmt.Println("main goroutine")
	time.Sleep(1 * time.Second)
}
