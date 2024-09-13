package main

import "fmt"

func add(a, b int, c ...int) (sum int, err error) {
	sum = a + b

	for _, v := range c {
		sum += v
	}

	return
}

func autoIncrement() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {

	a, b := 1, 2
	fmt.Println(add(a, b, 4, 5, 6))

	// 自动递增
	inc := autoIncrement()
	fmt.Println(inc())
	fmt.Println(inc())
	fmt.Println(inc())
}
