package data

import (
	"context"
	"shop/app/shop/api/internal_srv/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type CategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.CategoryBrandDOList, error)
	GetBrandList(ctx context.Context, categoryID int64) (*do.CategoryBrandDOList, error)
	Get(ctx context.Context, id int64) (*do.CategoryBrandDO, error)
	Create(ctx context.Context, gcb *do.CategoryBrandDO) error
	Update(ctx context.Context, gcb *do.CategoryBrandDO) error
	Delete(ctx context.Context, ID uint64) error
}
