package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint32, orderby []string) (*do.GoodsDOList, error)
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	Delete(ctx context.Context, ID uint64) error
}
