package db

import (
	"context"
	"shop/app/shop_srv/inventory/srv/internal/data/v1"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"

	"gorm.io/gorm"
)

type inventory struct {
	db *gorm.DB
}

func (i *inventory) Get(ctx context.Context, goodsID int64) (*do.InventoryDO, error) {
	db := i.db.WithContext(ctx)
	var inv do.InventoryDO

	err := db.Where("goods = ?", goodsID).First(&inv).Error
	if err != nil {
		//log.Errorf("get inv err: %v", err)
		return nil, err
	}

	return &inv, nil
}

func (i *inventory) Create(ctx context.Context, txn *gorm.DB, inventoryDO *do.InventoryDO) *gorm.DB {
	db := i.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(&inventoryDO)
}

func (i *inventory) GetSellDetail(ctx context.Context, ordersn string) (*do.StockSellDetailDO, error) {
	db := i.db.WithContext(ctx)
	var ordersellDetail do.StockSellDetailDO

	err := db.Where("order_sn = ?", ordersn).First(&ordersellDetail).Error
	if err != nil {
		return nil, err
	}
	return &ordersellDetail, err
}

func (i *inventory) Reduce(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) *gorm.DB {
	db := i.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).Where("stocks >= ?", num).UpdateColumn("stocks", gorm.Expr("stocks - ?", num))
}

func (i *inventory) Increase(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) *gorm.DB {
	db := i.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).UpdateColumn("stocks", gorm.Expr("stocks + ?", num))
}

func (i *inventory) CreateStockSellDetail(ctx context.Context, txn *gorm.DB, detail *do.StockSellDetailDO) *gorm.DB {
	db := i.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(&detail)
}

func (i *inventory) UpdateStockSellDetailStatus(ctx context.Context, txn *gorm.DB, ordersn string, status int32) *gorm.DB {
	db := i.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Model(do.StockSellDetailDO{}).Where("order_sn = ?", ordersn).Update("status", status)
}

func newInventory(factory *mysqlFactory) data.InventoryStore {
	return &inventory{db: factory.db}
}
