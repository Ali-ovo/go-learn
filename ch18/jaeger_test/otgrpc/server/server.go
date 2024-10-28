package main

import (
	"context"
	"go-learn/ch14/grpc/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	time.Sleep(2 * time.Second)

	return &proto.HelloReply{
		Message: "Hello " + request.Name,
		Data: &proto.HelloReply_Result{
			Name: "Ali",
			Url:  "https:/test.com",
		},
	}, nil
}

func main() {
	g := grpc.NewServer()

	proto.RegisterGreeterServer(g, &Server{})

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	err = g.Serve(listen)
	if err != nil {
		panic(err)
	}
}
