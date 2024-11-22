package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type GoodsCategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandDOList, error)
	Create(ctx context.Context, gcb *do.GoodsCategoryBrandDOList) error
	Update(ctx context.Context, gcb *do.GoodsCategoryBrandDOList) error
	Delete(ctx context.Context, ID uint64) error
}
