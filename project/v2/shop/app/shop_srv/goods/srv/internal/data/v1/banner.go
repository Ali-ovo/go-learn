package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
)

type BannerStore interface {
	List(ctx context.Context) (*do.BannerDOList, error)
	Get(ctx context.Context, id int64) (*do.BannerDO, error)
	Create(ctx context.Context, banner *do.BannerDO) error
	Update(ctx context.Context, banner *do.BannerDO) error
	Delete(ctx context.Context, ID int64) error
}
