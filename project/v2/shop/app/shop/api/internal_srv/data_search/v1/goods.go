package data_search

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_srv/domain/do"
)

type GoodsFilterRequest struct {
	*goods_pb.GoodsFilterRequest
	CategoryIDs []any
}

type GoodsStore interface {
	Search(ctx context.Context, request *GoodsFilterRequest) (*do.GoodsSearchDOList, error)
	Create(ctx context.Context, goods *do.GoodsSearchDO) error
	Update(ctx context.Context, goods *do.GoodsSearchDO) error
	Delete(ctx context.Context, ID uint64) error
}
