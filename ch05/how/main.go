package main

import "fmt"

func printSlice(data []string) {
	data[0] = "java"
}

func main() {
	// go 的切片在函数传参是值传递还是引用传递
	courses := []string{"go", "python", "java"}
	printSlice(courses)
	fmt.Println(courses)

	data := []int{1, 2, 3, 4, 5}
	s1 := data[1:3]
	s2 := data[1:]
	s2[0] = 22

	fmt.Println(s1) // [22 3]
	fmt.Println(s2) // [22 3 4 5]

	s2 = append(s2, 11, 22, 33, 44, 55, 66, 77)
	s2[0] = 222

	fmt.Println(s1) // [22 3]
	fmt.Println(s2) // [222 3 4 5 11 22 33 44 55 66 77]

}
