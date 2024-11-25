package goods

import (
	"shop/app/shop_api/api/internal/data/v2/goodsSrv"
	Igoods "shop/app/shop_api/api/internal/data/v2/goodsSrv/goods"
	"shop/app/shop_api/api/internal/data/v2/rpc"
)

type goodsFactory struct {
	gf rpc.GrpcFactory
}

func (g goodsFactory) Goods() Igoods.Goods {
	return newGoods(g.gf)
}

func (g goodsFactory) Brand() Igoods.Brand {
	return newBrand(g.gf)
}

func NewGoods(gf rpc.GrpcFactory) goodsSrv.GoodsData {
	return goodsFactory{gf}
}
