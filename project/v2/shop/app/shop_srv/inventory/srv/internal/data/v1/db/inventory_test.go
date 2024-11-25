package db

import (
	"context"
	"fmt"
	"shop/app/shop_srv/inventory/srv/internal/data/v1"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/pkg/options"
	"testing"
	"time"
)

func Conn() data.DataFactory {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.16.192",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "shop_inventory_srv",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              4,
		EnableLog:             true,
	}
	dbFactory, err := GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}
	return dbFactory
}

func TestInventory_Reduce(t *testing.T) {
	dbFactory := Conn()
	tx := dbFactory.Begin()
	goodID := 421
	goodNum := 2
	// UPDATE `inventory` SET `stocks`=stocks - 2 WHERE goods=421 AND stocks >= 2 AND `inventory`.`deleted_at` IS NULL
	err := dbFactory.Inventory().Reduce(context.Background(), tx, int64(goodID), int32(goodNum))
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
	tx.Commit()
	fmt.Printf("执行 库存扣减成功  %d : %d\n", goodID, goodNum)
	time.Sleep(time.Second)
}

func TestInventory_Increase(t *testing.T) {
	dbFactory := Conn()
	tx := dbFactory.Begin()
	goodID := 421
	goodNum := 2
	// UPDATE `inventory` SET `stocks`=stocks + 2 WHERE goods=421 AND `inventory`.`deleted_at` IS NULL
	err := dbFactory.Inventory().Increase(context.Background(), tx, int64(goodID), int32(goodNum))
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
	tx.Commit()
	fmt.Printf("执行 库存新增成功  %d : %d\n", goodID, goodNum)
}

func TestInventory_UpdateStockSellDetailStatus(t *testing.T) {
	dbFactory := Conn()
	tx := dbFactory.Begin()
	orderSn := "111"
	orderStatus := 1
	// UPDATE `stockselldetail` SET `status`=1 WHERE order_sn = '111'
	err := dbFactory.Inventory().UpdateStockSellDetailStatus(context.Background(), tx, orderSn, int32(orderStatus))
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}

	tx.Commit()
	fmt.Printf("执行 订单状态修改成功  %s : %d\n", orderSn, orderStatus)
}

func TestInventory_CreateStockSellDetail(t *testing.T) {
	dbFactory := Conn()
	tx := dbFactory.Begin()

	sellDetail := do.StockSellDetailDO{
		OrderSn: "111",
		Status:  1,
		Detail: do.GoodsDetailList{
			{Goods: 421, Num: 2},
			{Goods: 422, Num: 3},
		},
	}

	// INSERT INTO `stockselldetail` (`order_sn`,`status`,`detail`) VALUES ('111',1,'[{"Goods":421,"Num":2},{"Goods":422,"Num":3}]')
	err := dbFactory.Inventory().CreateStockSellDetail(context.Background(), tx, &sellDetail)
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}

	tx.Commit()
	fmt.Printf("执行 订单创建成功\n")
}
