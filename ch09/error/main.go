package main

import "fmt"

func mPrint(datas ...interface{}) {
	for _, value := range datas {
		fmt.Println(value)
	}

}

func main() {
	var data = []interface{}{
		"bobby", 18, 1.80,
	}

	mPrint(data...)

}
