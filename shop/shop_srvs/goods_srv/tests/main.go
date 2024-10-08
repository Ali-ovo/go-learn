package main

import (
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

func main() {
	Init()

	// TestGetBrandList()
	// TestCategoryList()
	TestGetSubCategory()

	defer conn.Close()
}
