package data

import (
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_api/data/v1/userSrv"
)

type DataFactory interface {
	User() userSrv.UserData
	Goods() goods_pb.GoodsClient
	//Goods() goodsSrv.GoodsData
}
