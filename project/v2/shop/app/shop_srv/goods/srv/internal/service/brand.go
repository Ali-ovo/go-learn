package service

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
)

type BrandSrv interface {
	List(ctx context.Context, request *goods_pb.BrandFilterRequest) (*dto.BrandsDTOList, error)
	Create(ctx context.Context, brand *dto.BrandsDTO) (int64, error)
	Update(ctx context.Context, brand *dto.BrandsDTO) error
	Delete(ctx context.Context, id int64) error
}
