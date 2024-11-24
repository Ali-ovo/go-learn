package service

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
)

type CategorySrv interface {
	AllList(ctx context.Context) (*dto.CategoryDTOList, error)
	List(ctx context.Context, level int32) (*dto.CategoryDTOList, error)
	Get(ctx context.Context, id int64) (*dto.CategoryDTO, error)
	Create(ctx context.Context, dto *dto.CategoryDTO) (int64, error)
	Update(ctx context.Context, dto *dto.CategoryDTO) error
	Delete(ctx context.Context, id int64) error
}
