package main

import "fmt"

func main() {
	name := "go学习"
	bytes := []rune(name)
	fmt.Println(len(bytes))

	// 转义
	courseName := "hello\r\ngo课程"
	fmt.Println(courseName)
	fmt.Println("hello\r\n")
	fmt.Println("world\r\n")

	// 格式化输出
	username := "ali"
	age := 18
	address := "上海"

	fmt.Printf("用户名：%s, 年龄： %d, 地址: %s \r\n", username, age, address)

	userMsg := fmt.Sprintf("用户名：%v, 年龄： %d, 地址: %s \r\n", username, age, address)

	fmt.Print(userMsg)

}
