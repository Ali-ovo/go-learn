package main

import "fmt"

func main() {

	// switch 语句
	a := 1
	switch a {
	case 1:
		fmt.Println("1")

	default:
		fmt.Println("other")
	}

	source := 110

	switch {
	case source < 60:
		fmt.Println("D")
	case source < 80:
		fmt.Println("C")
	case source < 90:
		fmt.Println("B")
	case source < 100:
		fmt.Println("A")
	default:
		fmt.Println("S")
	}
}
