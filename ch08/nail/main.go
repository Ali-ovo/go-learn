package main

import "fmt"

type Person struct {
	name string
	age  int
}

func main() {

	p1 := Person{
		name: "zhangsan",
		age:  18,
	}

	p2 := Person{
		name: "zhangsan",
		age:  18,
	}

	fmt.Println(p1 == p2)

}
