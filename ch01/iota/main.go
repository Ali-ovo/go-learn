package main

import "fmt"

func main() {
	const (
		// iota 特殊常量 可以被编辑器修改 自动递增
		a = iota + 1
		b
		c = 4
		d = iota
	)
	fmt.Println(a, b, c, d)

	//  _ 代表匿名变量
	const (
		e = iota
		f
		_
		g
	)
	fmt.Println(e, f, g)

}
