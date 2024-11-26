package db

import (
	"context"
	"shop/app/shop_srv/inventory/srv/internal/data/v1"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type inventory struct {
	db *gorm.DB
}

func (i *inventory) Get(ctx context.Context, goodsID int64) (*do.InventoryDO, error) {
	inv := do.InventoryDO{}
	err := i.db.Where("goods = ?", goodsID).First(&inv).Error
	if err != nil {
		//log.Errorf("get inv err: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrInventoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return &inv, nil
}

func (i *inventory) Create(ctx context.Context, inventoryDO *do.InventoryDO) error {
	//设置库存， 如果我要更新库存
	tx := i.db.Create(&inventoryDO)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (i *inventory) GetSellDetail(ctx context.Context, txn *gorm.DB, ordersn string) (*do.StockSellDetailDO, error) {
	db := i.db
	if txn != nil {
		db = txn
	}
	var ordersellDetail do.StockSellDetailDO
	err := db.Where("order_sn = ?", ordersn).First(&ordersellDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrInvSellDetailNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &ordersellDetail, err
}

func (i *inventory) Reduce(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) *gorm.DB {
	db := i.db
	if txn != nil {
		db = txn
	}
	return db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).Where("stocks >= ?", num).UpdateColumn("stocks", gorm.Expr("stocks - ?", num))
}

func (i *inventory) Increase(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) error {
	db := i.db
	if txn != nil {
		db = txn
	}
	err := db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).UpdateColumn("stocks", gorm.Expr("stocks + ?", num)).Error
	return err
}

func (i *inventory) CreateStockSellDetail(ctx context.Context, txn *gorm.DB, detail *do.StockSellDetailDO) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	tx := db.Create(&detail)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (i *inventory) UpdateStockSellDetailStatus(ctx context.Context, txn *gorm.DB, ordersn string, status int32) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	//update语句如果没有更新的话那么不会报错，但是他会返回一个影响的行数，所以我们可以根据影响的行数来判断是否更新成功
	result := db.Model(do.StockSellDetailDO{}).Where("order_sn = ?", ordersn).Update("status", status)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}

	//这里应该在service层去写代码判断更合理
	//有两种情况都会导致影响的行数为0，一种是没有找到，一种是没有更新
	//if result.RowsAffected == 0 {
	//	return errors.WithCode(code.ErrInvSellDetailNotFound, "inventory sell detail not found")
	//}
	return nil
}

func newInventory(factory *mysqlFactory) data.InventoryStore {
	return &inventory{db: factory.db}
}
