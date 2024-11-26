package data

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type ShopCartStore interface {
	List(ctx context.Context, userID int64, opts metav1.ListMeta, orderby []string) ([]*do.ShoppingCartDOList, error)
	Create(ctx context.Context, cartItem *do.ShoppingCartDO) error
	Get(ctx context.Context, userID, goodsID int64) (*do.ShoppingCartDO, error)
	UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error
	Delete(ctx context.Context, userID int64) error
	ClearCheck(ctx context.Context, userID int64) error
}
