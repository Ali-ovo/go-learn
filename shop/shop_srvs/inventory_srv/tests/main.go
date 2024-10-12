package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/inventory_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var inventoryClient proto.InventoryClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:55137", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	inventoryClient = proto.NewInventoryClient(conn)
}

func TestSetInv(GoodsId int32, Num int32) {
	rsp, err := inventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: GoodsId,
		Num:     Num,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp, "SetInv")
}

func TestInvDetail() {
	rsp, err := inventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 1,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp, "InvDetail")
}

func TestSell() {
	rsp, err := inventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 2,
				Num:     10,
			},
			{
				GoodsId: 1,
				Num:     10,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp, "Sell")
}

func TestReback() {
	rsp, err := inventoryClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 2,
				Num:     10,
			},
			{
				GoodsId: 1,
				Num:     10,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp, "Reback")
}

func main() {
	Init()

	// TestSetInv()
	// TestInvDetail()
	// TestSell()
	// TestReback()

	for i := 421; i <= 840; i++ {
		TestSetInv(int32(i), 100)
	}

	defer conn.Close()
}
