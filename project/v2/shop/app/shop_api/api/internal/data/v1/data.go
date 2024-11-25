package data

import (
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_api/api/internal/data/v1/userSrv"
)

type DataFactory interface {
	User() userSrv.UserData
	Goods() goods_pb.GoodsClient
	//Goods() goodsSrv.GoodsData
}
