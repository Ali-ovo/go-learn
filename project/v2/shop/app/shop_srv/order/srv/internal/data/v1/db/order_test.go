package db

import (
	"context"
	"fmt"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	"shop/pkg/gorm"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestOrder_Create(t *testing.T) {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.16.192",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "shop_order_srv",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              4,
		EnableLog:             true,
	}
	dbConn, err := NewOrderSQLClient(mysqlOpts)
	factory := dataFactory{db: dbConn}
	order := newOrders(&factory)

	if err != nil {
		panic(err)
	}

	time := time.Now()

	txn := dbConn.Begin()
	//sourceTx := txn.Statement.ConnPool.(*sql.Tx)

	db := order.Create(context.Background(), txn, &do.OrderInfoDO{
		BaseModel: gorm.BaseModel{ID: 2},
		OrderGoods: []*do.OrderGoods{
			{
				Order:      2,
				Goods:      421,
				GoodsName:  "测试",
				GoodsImage: "测试",
				GoodsPrice: 11.111,
				Nums:       2,
			},
		},
		User:         3,
		OrderSn:      "111111",
		PayType:      "alipay",
		Status:       "PAYING",
		TradeNo:      "222222",
		OrderMount:   22.222,
		PayTime:      &time,
		Address:      "测试",
		SignerName:   "测试",
		SingerMobile: "13067353692",
		Post:         "测试",
	})
	if db.Error != nil {
		t.Error(db.Error)
	}
	txn.Commit()
	fmt.Println(db)
}
