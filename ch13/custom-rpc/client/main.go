package main

import (
	"fmt"
	"go-learn/ch13/custom-rpc/client_proxy"
)

func main() {
	// 建立连接
	// client, err := rpc.Dial("tcp", ":1234")
	client := client_proxy.NewHelloServiceClient("tcp", ":1234")

	var replay string
	// err = client.Call(handler.HelloServiceName+".Hello", "world", &replay)

	err := client.Hello("world", &replay)

	if err != nil {
		panic(err)
	}
	fmt.Println(replay)
}
