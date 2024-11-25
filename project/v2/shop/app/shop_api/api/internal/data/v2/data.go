package data

import (
	"shop/app/shop_api/api/internal/data/v2/goodsSrv"
	"shop/app/shop_api/api/internal/data/v2/userSrv"
)

type DataFactory interface {
	User() userSrv.UserData
	Goods() goodsSrv.GoodsData
	//Goods() goodsSrv.GoodsData
}
