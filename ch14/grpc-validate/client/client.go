package main

import (
	"context"
	"fmt"
	"go-learn/ch14/grpc-validate/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type customCredential struct{}

func main() {
	var opts []grpc.DialOption

	ops := append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:50051", ops...)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	//rsp, _ := c.Search(context.Background(), &empty.Empty{})
	rsp, err := c.SayHello(context.Background(), &proto.Person{
		Id:     1000,
		Email:  "ali_ovo@qq.com",
		Mobile: "18888888888",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Email)
}
