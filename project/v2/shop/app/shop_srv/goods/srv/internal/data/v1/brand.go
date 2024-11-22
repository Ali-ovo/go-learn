package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type BrandsStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error)
	Create(ctx context.Context, brands *do.BrandsDO) error
	Update(ctx context.Context, brands *do.BrandsDO) error
	Delete(ctx context.Context, ID uint64) error
}
