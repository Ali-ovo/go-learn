package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
)

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	Delete(ctx context.Context, ID uint64) error
}
