package data

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type OrderStore interface {
	Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error)
	List(ctx context.Context, userID int64, meta metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error)
	Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) *gorm.DB
	Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) *gorm.DB
}
