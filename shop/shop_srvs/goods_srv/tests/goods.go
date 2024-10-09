package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/goods_srv/proto"
)

func TesGoodsList() {
	rsp, err := brandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		// KeyWords:    "深海速冻",
		PriceMin: 90,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	for _, v := range rsp.Data {
		fmt.Println(v.Name, v.ShopPrice)
	}
}

func TesBatchGetGoods() {
	rsp, err := brandClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	for _, v := range rsp.Data {
		fmt.Println(v.Name, v.ShopPrice)
	}
}

func TesGetGoodsDetail() {
	rsp, err := brandClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name)

}
