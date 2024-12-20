package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"

	"gorm.io/gorm"
)

type BannerStore interface {
	List(ctx context.Context) (*do.BannerDOList, error)
	Get(ctx context.Context, id int64) (*do.BannerDO, error)
	Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) *gorm.DB
	Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) *gorm.DB
	Delete(ctx context.Context, txn *gorm.DB, ID int64) *gorm.DB
}
