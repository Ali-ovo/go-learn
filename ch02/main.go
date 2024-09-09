package main

import "fmt"

func main() {
	var a int8 = 127
	var b int16
	var c int32
	var d int64

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)

	//  uint 无符号整数 不会为负数 并且 >= 0
	var ua uint8
	var ub uint16
	var uc uint32
	var ud uint64

	fmt.Println(ua)
	fmt.Println(ub)
	fmt.Println(uc)
	fmt.Println(ud)

	var f1 float32
	var f2 float64
	fmt.Println(f1)
	fmt.Println(f2)

	// byte 本质是 uint8  用来存放字符
	var bytec byte = 'c'
	fmt.Println(bytec) // 99

	// rune 本质是 int32  用来存放 中文或者去其他国家的字符
	var runea rune = 'a'
	fmt.Println(runea) // 97
}
