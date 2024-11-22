package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
)

type CategoryStore interface {
	Get(ctx context.Context, ID int32) (*do.CategoryDO, error)
	ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error)
	Create(ctx context.Context, category *do.CategoryDO) error
	Update(ctx context.Context, category *do.CategoryDO) error
	Delete(ctx context.Context, ID int64) error
}
