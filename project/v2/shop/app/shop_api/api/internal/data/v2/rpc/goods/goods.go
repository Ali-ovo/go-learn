package goods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	Igoods "shop/app/shop_api/api/internal/data/v2/goodsSrv/goods"
	"shop/app/shop_api/api/internal/data/v2/rpc"
	doGoods "shop/app/shop_api/api/internal/domain/do/goods"
)

type goods struct {
	gp goods_pb.GoodsClient
}

func (g goods) List(ctx context.Context, request *goods_pb.GoodsFilterRequest) (*doGoods.GoodDOList, error) {
	//TODO implement me
	panic("implement me")
}

func (g goods) GetBatch(ctx context.Context, ids []int64) (*doGoods.GoodDOList, error) {
	//TODO implement me
	panic("implement me")
}

func (g goods) Get(ctx context.Context, id int64) (*doGoods.GoodDO, error) {
	//TODO implement me
	panic("implement me")
}

func (g goods) Create(ctx context.Context, request *doGoods.GoodDO) (*doGoods.GoodDO, error) {
	//TODO implement me
	panic("implement me")
}

func (g goods) Update(ctx context.Context, request *doGoods.GoodDO) error {
	//TODO implement me
	panic("implement me")
}

func (g goods) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func newGoods(gf rpc.GrpcFactory) Igoods.Goods {
	return &goods{
		gp: gf.GoodsClient,
	}
}
