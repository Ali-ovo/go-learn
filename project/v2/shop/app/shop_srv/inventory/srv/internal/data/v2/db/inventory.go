package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"shop/app/shop_srv/inventory/srv/internal/data/v2"
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

func (i *inventory) GetSellDetail(ctx context.Context, ordersn string) (*do.StockSellDetailDO, error) {
	var ordersellDetail do.StockSellDetailDO
	err := i.db.Where("order_sn = ?", ordersn).First(&ordersellDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrInvSellDetailNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &ordersellDetail, err
}

func (i *inventory) Reduce(ctx context.Context, txn *sql.Tx, goodsID int64, num int32) (sql.Result, error) {
	result, err := txn.ExecContext(ctx, "UPDATE `inventory` SET `stocks`=stocks - ? WHERE goods=? AND stocks >= ? AND `inventory`.`deleted_at` IS NULL", num, goodsID, num)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return result, nil
}

func (i *inventory) Increase(ctx context.Context, txn *sql.Tx, goodsID int64, num int32) (sql.Result, error) {
	result, err := txn.ExecContext(ctx, "UPDATE `inventory` SET `stocks`=stocks + ? WHERE goods=? AND `inventory`.`deleted_at` IS NULL", num, goodsID)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return result, nil
}

func (i *inventory) CreateStockSellDetail(ctx context.Context, txn *sql.Tx, detail *do.StockSellDetailDO) (sql.Result, error) {
	// INSERT INTO `stockselldetail` (`order_sn`,`status`,`detail`) VALUES ('111',1,'[{"Goods":421,"Num":2},{"Goods":422,"Num":3}]')
	detailByte, err := json.Marshal(detail.Detail)
	if err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}

	result, err := txn.ExecContext(ctx, "INSERT INTO `stockselldetail` (`order_sn`,`status`,`detail`) VALUES (?,?,?)", detail.OrderSn, detail.Status, string(detailByte))
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return result, nil
}

func (i *inventory) UpdateStockSellDetailStatus(ctx context.Context, txn *sql.Tx, ordersn string, status int32) (sql.Result, error) {
	//update语句如果没有更新的话那么不会报错，但是他会返回一个影响的行数，所以我们可以根据影响的行数来判断是否更新成功
	result, err := txn.ExecContext(ctx, "UPDATE `stockselldetail` SET `status`=? WHERE order_sn = ?", status, ordersn)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return result, nil
}

func newInventory(factory *mysqlFactory) data.InventoryStore {
	return &inventory{db: factory.db}
}
