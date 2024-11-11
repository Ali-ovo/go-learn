package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "shop/api/user/v1"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	uc := v1.NewUserClient(conn)
	list, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}
