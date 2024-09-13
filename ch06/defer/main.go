package main

import (
	"fmt"
)

func main() {
	// defer

	// var mu sync.Mutex
	// mu.Lock()
	// defer mu.Unlock()

	defer fmt.Println("1")
	defer fmt.Println("2")
	fmt.Println("3")

}
