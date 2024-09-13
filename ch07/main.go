package main

import (
	"fmt"
	"strconv"
)

// 自定义类型
type Myint3 int

func (mi Myint3) string() string {
	return strconv.Itoa(int(mi))
}

func main() {

	// type 关键字
	type MyInt = int
	var i MyInt
	fmt.Println(i)

	// 自定义类型
	type MyInt2 int
	var j MyInt2
	fmt.Printf("MyInt2: %T\r\n", j)

	var k Myint3
	fmt.Printf("MyInt3: %T\r\n", k.string())
}
