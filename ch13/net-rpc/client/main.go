package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 建立连接
	conn, err := rpc.Dial("tcp", "127.0.0.1:1234")

	if err != nil {
		panic(err)
	}

	var replay *string = new(string)
	// 调用远程方法
	err = conn.Call("HelloService.Hello", "world", replay)

	if err != nil {
		panic(err)
	}
	fmt.Println(*replay)
}
