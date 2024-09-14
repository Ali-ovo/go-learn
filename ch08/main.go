package main

import "fmt"

type Person struct {
	name string
}

func changeName(p *Person) {
	p.name = "lisi"
}

func swap(a, b *int) {
	*a, *b = *b, *a
}

func main() {

	p := Person{name: "zhangsan"}
	changeName(&p)

	fmt.Println(p)

	a, b := 1, 2
	swap(&a, &b)
	fmt.Println(a, b)
}
