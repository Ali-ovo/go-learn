package service

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
)

type CategoryBrandSrv interface {
	List(ctx context.Context, request *goods_pb.CategoryBrandFilterRequest) (*dto.CategoryBrandDTOList, error)
	Get(ctx context.Context, categoryID int64) (*dto.CategoryBrandDTOList, error)
	Create(ctx context.Context, dto *dto.CategoryBrandDTO) (int64, error)
	Update(ctx context.Context, dto *dto.CategoryBrandDTO) error
	Delete(ctx context.Context, id int64) error
}
