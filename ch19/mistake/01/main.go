package main

import "fmt"

func main() {

	// map 必须初始化
	// var course map[string]string

	var course = make(map[string]string, 2)
	course["name"] = "Golang"
	course["time"] = "2020-01-01"

	fmt.Println(course["desc"])

}
