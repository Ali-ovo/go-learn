package db

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type orders struct {
	db *gorm.DB
}

func (o *orders) Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error) {
	db := o.db.WithContext(ctx)
	var order do.OrderInfoDO

	err := db.Preload("OrderGoods").Where("order_sn = ?", orderSn).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *orders) List(ctx context.Context, userID int64, opts metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error) {
	db := o.db.WithContext(ctx)
	var ret do.OrderInfoDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.OrderInfoDO{})
	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count

	return &ret, nil
}

// Create
//
//	@Description: 创建订单 (需要删除对应的购物车记录)
//	@receiver o
//	@param ctx
//	@param txn
//	@param order
//	@return error
func (o *orders) Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) *gorm.DB {
	db := o.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(order)
}

func (o *orders) Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) *gorm.DB {
	db := o.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	return db.Where("order_sn", order.OrderSn).Updates(order)
}

func newOrders(factory *dataFactory) data.OrderStore {
	return &orders{
		db: factory.db,
	}
}
