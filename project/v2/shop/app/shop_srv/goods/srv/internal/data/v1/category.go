package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"

	"gorm.io/gorm"
)

type CategoryStore interface {
	Get(ctx context.Context, ID int64) (*do.CategoryDO, error)
	List(ctx context.Context, level int32) (*do.CategoryDOList, error)
	ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error)
	Create(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) error
	Update(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) error
	Delete(ctx context.Context, txn *gorm.DB, ID int64) error
}
