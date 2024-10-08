package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/goods_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	brandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()

	TestGetBrandList()

	defer conn.Close()
}
