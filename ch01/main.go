package main

import "fmt"

// 全局变量
const (
	PI = 3.14
	PI1
	name = "zhangsan"
	name2
)

func main() {
	age := 1

	var user1, user2 = "zhangsan", "lisi"

	fmt.Println(age)
	fmt.Println(PI, name)
	fmt.Println(user1, user2)

	// same PI name
	fmt.Println(PI1, name2)

	// add test2
}
