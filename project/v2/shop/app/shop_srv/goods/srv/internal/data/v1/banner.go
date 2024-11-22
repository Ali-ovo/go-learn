package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type BannerStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerDOList, error)
	Create(ctx context.Context, banner *do.BannerDO) error
	Update(ctx context.Context, banner *do.BannerDO) error
	Delete(ctx context.Context, ID int64) error
}
