package goods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	Igoods "shop/app/shop_api/api/internal/data/v2/goodsSrv/goods"
	"shop/app/shop_api/api/internal/data/v2/rpc"
	doGoods "shop/app/shop_api/api/internal/domain/do/goods"
)

type brand struct {
	gp goods_pb.GoodsClient
}

func (b brand) List(ctx context.Context, in *goods_pb.BrandFilterRequest) (*doGoods.BrandDOList, error) {
	//TODO implement me
	panic("implement me")
}

func (b brand) Create(ctx context.Context, request *doGoods.BrandDO) (*doGoods.BrandDO, error) {
	//TODO implement me
	panic("implement me")
}

func (b brand) Update(ctx context.Context, request *doGoods.BrandDO) error {
	//TODO implement me
	panic("implement me")
}

func (b brand) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func newBrand(gf rpc.GrpcFactory) Igoods.Brand {
	return &brand{
		gp: gf.GoodsClient,
	}
}
