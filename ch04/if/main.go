package main

import "fmt"

func main() {
	// if bool
	a := 18
	country := "china"
	if a < 18 && country == "china" {
		fmt.Println("未成年")
	} else {
		fmt.Println("成年")
	}

}
