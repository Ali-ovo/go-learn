package db

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type shopCarts struct {
	db *gorm.DB
}

func (sc *shopCarts) List(ctx context.Context, userID int64, checked bool, opts metav1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error) {
	db := sc.db.WithContext(ctx)
	var ret do.ShoppingCartDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.ShoppingCartDO{}).Where("user = ?", userID)
	if checked {
		result = result.Where("checked = ?", true)
	}

	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count

	return &ret, nil
}

func (sc *shopCarts) Create(ctx context.Context, txn *gorm.DB, cartItem *do.ShoppingCartDO) *gorm.DB {
	db := sc.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(cartItem)
}

func (sc *shopCarts) Get(ctx context.Context, userID, goodsID int64) (*do.ShoppingCartDO, error) {
	db := sc.db.WithContext(ctx)
	var shopCart do.ShoppingCartDO

	result := db.Where("user = ? AND goods = ?", userID, goodsID).First(&shopCart)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shopCart, nil
}

func (sc *shopCarts) Update(ctx context.Context, txn *gorm.DB, cartItem *do.ShoppingCartDO) *gorm.DB {
	db := sc.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Model(do.ShoppingCartDO{}).Where("user = ? AND goods = ?", cartItem.User, cartItem.Goods).Updates(do.ShoppingCartDO{Nums: cartItem.Nums, Checked: cartItem.Checked})
}

func (sc *shopCarts) Delete(ctx context.Context, txn *gorm.DB, userID int64) *gorm.DB {
	db := sc.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Where("user = ?", userID).Delete(&do.ShoppingCartDO{})
}

// DeleteByGoodsIDs 删除选中商品的购物车记录, 下订单
// 2 种做法
// 下单后, 直接执行删除购物车的记录
// 下单后什么都不做, 直接给 rocketmq 发送一个消息, 然后由 rocketmq 来执行删除购物车的记录 (因为使用 dtm 我就没有使用 rocketmq)
func (sc *shopCarts) DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID int64, goodsIDs []int64) *gorm.DB {
	db := sc.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Where("user = ? AND goods IN (?)", userID, goodsIDs).Delete(&do.ShoppingCartDO{})
}

var _ data.ShopCartStore = &shopCarts{}

func newShopCarts(factory *dataFactory) data.ShopCartStore {
	return &shopCarts{
		db: factory.db,
	}
}
