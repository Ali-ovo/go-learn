package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:1234")

	if err != nil {
		panic(err)
	}

	var replay *string = new(string)
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err = client.Call("HelloService.Hello", "world", replay)

	if err != nil {
		panic(err)
	}
	fmt.Println(*replay)
}
