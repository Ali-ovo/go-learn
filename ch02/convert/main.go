package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 强制转换
	var a int = 1
	var b int64 = int64(a)
	var c int32 = int32(b)

	fmt.Println(a, b, c)

	var f1 float32 = 3.14
	var f2 int32 = int32(f1)
	var f3 float64 = float64(f2)

	fmt.Println(f1, f2, f3)

	var s1 string = "12"
	s2, _ := strconv.Atoi(s1)
	fmt.Println(s1, s2)

	var int1 int = 99
	sint1 := strconv.Itoa(int1)
	fmt.Println(int1, sint1)

	// string to float
	s4, _ := strconv.ParseFloat("3.14", 32)
	fmt.Println(s4)

	// string to int
	s5, _ := strconv.ParseInt("-12", 10, 64)
	fmt.Println(s5)

	// string to bool
	b1, _ := strconv.ParseBool("true")
	fmt.Println(b1)

	// 基本类型转字符串
	boolStr := strconv.FormatBool(true)
	fmt.Println(boolStr)

	flootStr := strconv.FormatFloat(3.14, 'f', 2, 64)
	fmt.Println(flootStr)

	intStr := strconv.FormatInt(-12, 10)
	fmt.Println(intStr)
}
