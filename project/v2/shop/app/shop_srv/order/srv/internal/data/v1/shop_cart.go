package data

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type ShopCartStore interface {
	List(ctx context.Context, userID int64, checked bool, opts metav1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error) // 查看用户购物车信息
	Create(ctx context.Context, txn *gorm.DB, cartItem *do.ShoppingCartDO) *gorm.DB                                               // 创建购物车商品
	Get(ctx context.Context, userID, goodsID int64) (*do.ShoppingCartDO, error)                                                   // 获取购物车商品详情
	Update(ctx context.Context, txn *gorm.DB, cartItem *do.ShoppingCartDO) *gorm.DB                                               // 更新购物车商品
	Delete(ctx context.Context, txn *gorm.DB, userID int64) *gorm.DB                                                              // 清空购物车
	DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID int64, goodsIDs []int64) *gorm.DB                                  // 删除指定 购物车商品
}
