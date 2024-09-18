package main

import "fmt"

type Company struct {
	Name    string
	Address string
}

type PrintResult struct {
	Info string
	Err  error
}

type Employee struct {
	Name    string
	company Company
}

func RpcPrintln(employee Employee) {

	fmt.Println(employee)
}

func main() {

}
