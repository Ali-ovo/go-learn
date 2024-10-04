package main

import (
	"context"
	"go-learn/ch18/grpclb_test/proto"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"consul://192.168.189.128:8500/user_srv?wait=14s&tag=srv",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userSrvClient := proto.NewUserClient(conn)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})

	if err != nil {
		panic(err)
	}

	for index, data := range rsp.Data {
		log.Printf("index: %d, data: %v", index, data)
	}

}
