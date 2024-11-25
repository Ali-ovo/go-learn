package data

import (
	"context"
	"shop/app/shop/api/internal_srv/domain/do"
)

type CategoryStore interface {
	Get(ctx context.Context, ID int64) (*do.CategoryDO, error)
	List(ctx context.Context, level int32) (*do.CategoryDOList, error)
	ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error)
	Create(ctx context.Context, category *do.CategoryDO) error
	Update(ctx context.Context, category *do.CategoryDO) error
	Delete(ctx context.Context, ID int64) error
}
