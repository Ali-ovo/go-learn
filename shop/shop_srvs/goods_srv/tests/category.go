package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	fmt.Println(rsp.JsonData)
}

func TestGetSubCategory() {
	rsp, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}
