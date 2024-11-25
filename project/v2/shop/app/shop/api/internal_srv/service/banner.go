package service

import (
	"context"
	"shop/app/shop/api/internal_srv/domain/dto"
)

type BannerSrv interface {
	List(ctx context.Context) (*dto.BannerDTOList, error)
	Create(ctx context.Context, branner *dto.BannerDTO) (int64, error)
	Update(ctx context.Context, branner *dto.BannerDTO) error
	Delete(ctx context.Context, id int64) error
}
