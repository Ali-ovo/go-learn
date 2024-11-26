package db

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"

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

func (o *orders) List(ctx context.Context, userID uint64, opts metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error) {
	db := o.db.WithContext(ctx)
	var ret do.OrderInfoDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	query := db.Model(&do.OrderInfoDO{})
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
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
func (o *orders) Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
	db := o.db.WithContext(ctx)
	if txn != nil {
		db = txn
	}

	db = db.Create(order)
	if db.Error != nil {
		return errors.WithCode(code.ErrDatabase, db.Error.Error())
	}
	return nil
}

func (o *orders) Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
	db := o.db.WithContext(ctx)
	if txn != nil {
		db = txn
	}

	db = db.Model(order).Updates(order)
	if db.Error != nil {
		return errors.WithCode(code.ErrDatabase, db.Error.Error())
	}
	return nil
}

var _ data.OrderStore = &orders{}
