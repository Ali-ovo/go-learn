package main

import "fmt"

type Person struct {
	name string
	age  int
}

type Student struct {
	// 默认嵌套
	person Person
	score  float32

	// 第二张嵌套
	Person
}

// 结构体方法
func (p Person) print() {
	fmt.Printf("Name: %s, Age: %d \r\n", p.name, p.age)
}

func main() {
	// 嵌套结构体
	s := Student{
		person: Person{
			name: "zhangsan",
			age:  18,
		},

		score: 100.0,
	}

	fmt.Println(s.age)

	p := Person{
		name: "zhangsan",
		age:  18,
	}
	p.print()
}
