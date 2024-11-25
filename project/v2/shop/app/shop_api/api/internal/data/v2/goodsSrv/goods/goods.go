package Igoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	doGoods "shop/app/shop_api/api/internal/domain/do/goods"
)

type Goods interface {
	List(ctx context.Context, request *goods_pb.GoodsFilterRequest) (*doGoods.GoodDOList, error)
	GetBatch(ctx context.Context, ids []int64) (*doGoods.GoodDOList, error)
	Get(ctx context.Context, id int64) (*doGoods.GoodDO, error)
	Create(ctx context.Context, request *doGoods.GoodDO) (*doGoods.GoodDO, error)
	Update(ctx context.Context, request *doGoods.GoodDO) error
	Delete(ctx context.Context, id int64) error
}
