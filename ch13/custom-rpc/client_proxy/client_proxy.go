package client_proxy

import (
	"go-learn/ch13/custom-rpc/handler"
	"net/rpc"
)

type HelloServiceStub struct {
	*rpc.Client
}

func NewHelloServiceClient(protol, address string) HelloServiceStub {
	conn, err := rpc.Dial(protol, address)
	if err != nil {
		panic(err)
	}

	return HelloServiceStub{Client: conn}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {

	err := c.Call(handler.HelloServiceName+".Hello", request, reply)

	if err != nil {
		return err
	}
	return nil
}
