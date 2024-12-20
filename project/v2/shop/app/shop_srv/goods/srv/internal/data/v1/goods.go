package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsDOList, error)
	Create(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) *gorm.DB
	Update(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) *gorm.DB
	Delete(ctx context.Context, txn *gorm.DB, ID uint64) *gorm.DB
}
