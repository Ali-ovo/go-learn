package main

import "fmt"

func main() {

	// goto 语句
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == 2 {
				goto breakHere
			}

			fmt.Println(i, j)
		}
	}

breakHere:
	fmt.Println("over")

}
