package main

import "fmt"

func Add[T int | float64](a, b T) T {
	return a + b
}

type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

type Man struct{}
type Woman struct{}
type Company[T Man | Woman] struct {
	Name string
	CEO  T
}

func main() {
	fmt.Println(Add(1, 2))
	fmt.Println(Add(1.2, 2.2))

	m := MyMap[int, float64]{1: 1.1, 2: 2.2}
	fmt.Println(m)

	company := Company[Man]{
		Name: "Google",
		CEO:  Man{},
	}
	fmt.Println(company)
}
