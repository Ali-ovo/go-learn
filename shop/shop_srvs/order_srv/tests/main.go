package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/order_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(rsp)
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})

	if err != nil {
		panic(err)
	}

	for _, item := range rsp.Data {
		fmt.Println(item)
	}
}

func TestUpdateCartItem(id int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: true,
	})

	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  21,
		Address: "ShangHai",
		Mobile:  "123456789",
		Post:    "尽快发货",
	})

	if err != nil {
		panic(err)
	}
}

func TestGetOrderDetail(orderId int32) {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.OrderInfo.OrderSn)

	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}
}

func TestOrderList() {
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 21,
	})

	if err != nil {
		panic(err)
	}

	for _, order := range rsp.Data {
		fmt.Println(order)
	}
}

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:61007", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	orderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()

	// TestCreateCartItem(21, 1, 421)
	// TestCreateCartItem(21, 1, 422)
	// TestCartItemList(21)
	// TestUpdateCartItem(2)
	// TestCreateOrder()
	// TestGetOrderDetail(1)

	TestOrderList()

	defer conn.Close()
}
