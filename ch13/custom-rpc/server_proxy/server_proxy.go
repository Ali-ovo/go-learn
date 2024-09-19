package server_proxy

import (
	"go-learn/ch13/custom-rpc/handler"
	"net/rpc"
)

type HelloService interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(srv HelloService) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}
