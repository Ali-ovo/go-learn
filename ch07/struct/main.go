package main

import "fmt"

type Person struct {
	name    string
	age     int
	address string
	height  float32
}

func main() {

	// 结构体
	person := Person{
		name:    "zhangsan",
		age:     18,
		address: "shanghai",
		height:  1.75,
	}

	fmt.Println(person)

	persons := []Person{
		{name: "zhangsan", age: 18, address: "shanghai", height: 1.75},
	}

	fmt.Println(persons)

	// 匿名结构体
	address := struct {
		province string
	}{
		province: "shanghai",
	}

	fmt.Println(address)

}
