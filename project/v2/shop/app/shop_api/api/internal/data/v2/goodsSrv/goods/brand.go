package Igoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	doGoods "shop/app/shop_api/api/internal/domain/do/goods"
)

type Brand interface {
	List(ctx context.Context, in *goods_pb.BrandFilterRequest) (*doGoods.BrandDOList, error)
	Create(ctx context.Context, request *doGoods.BrandDO) (*doGoods.BrandDO, error)
	Update(ctx context.Context, request *doGoods.BrandDO) error
	Delete(ctx context.Context, id int64) error
}
