package main

import "fmt"

func main() {

	// 数组
	var arr = [3]int{1, 2, 3}
	arr = [...]int{4, 5, 6}
	fmt.Println(arr)

	// 多维数组
	arr1 := [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	fmt.Println(arr1)

}
