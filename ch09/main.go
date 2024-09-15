package main

import "fmt"

type Duck interface {
	Gaga()
	Walk()
	Swimming()
}

type pskDuck struct {
	legs int
}

func (pd *pskDuck) Walk() {
	fmt.Println("Walk")

}

func (pd *pskDuck) Gaga() {
	fmt.Println("gaga")
}

func (pd *pskDuck) Swimming() {
	fmt.Println("swimming", pd.legs)
}

func main() {
	// duck类型

	var d Duck = &pskDuck{
		legs: 2,
	}
	d.Gaga()
	d.Walk()
	d.Swimming()

}
