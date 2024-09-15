package main

import (
	"fmt"
)

func add(a, b interface{}) interface{} {
	switch a.(type) {
	case int:
		ai, _ := a.(int)
		bi, _ := b.(int)
		return ai + bi

	case int32:
		ai, _ := a.(int32)
		bi, _ := b.(int32)
		return ai + bi
	}

	return nil
}

func main() {
	a, b := 1, 2

	fmt.Println(add(a, b))
}
