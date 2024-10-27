package main

import (
	"context"
	"fmt"
	"go-learn/ch14/grpc/proto"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	timestampFormat = time.StampNano // "Jan _2 15:04:05.000"
	streamingCount  = 10
)

func main() {

	retryOpts := []retry.CallOption{
		retry.WithMax(3),
		retry.WithPerRetryTimeout(1 * time.Second),
		retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}

	conn, err := grpc.NewClient(
		"127.0.0.1:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(retryOpts...)),
	)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)

	// md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	// ctx := metadata.NewOutgoingContext(context.Background(), md)

	r, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "ali",
		G:    proto.Gender_FEMALE,
		Mp: map[string]string{
			"key": "value",
		},
		AddTime: timestamppb.New(time.Now()),
	},
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(r.GetMessage())

}
