package main

import (
	"net"
	"net/rpc"

	"go-learn/ch13/custom-rpc/handler"
	"go-learn/ch13/custom-rpc/server_proxy"
)

func main() {
	listener, _ := net.Listen("tcp", ":1234")

	// _ = rpc.RegisterName(handler.HelloServiceName, &HelloService{})
	_ = server_proxy.RegisterHelloService(&handler.HelloService{})

	for {
		conn, _ := listener.Accept()
		go rpc.ServeConn(conn)
	}
}
