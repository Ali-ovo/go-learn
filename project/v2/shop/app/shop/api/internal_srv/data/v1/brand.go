package data

import (
	"context"
	"shop/app/shop/api/internal_srv/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type BrandsStore interface {
	Get(ctx context.Context, ID int64) (*do.BrandsDO, error)
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error)
	Create(ctx context.Context, brands *do.BrandsDO) error
	Update(ctx context.Context, brands *do.BrandsDO) error
	Delete(ctx context.Context, ID int64) error
}
