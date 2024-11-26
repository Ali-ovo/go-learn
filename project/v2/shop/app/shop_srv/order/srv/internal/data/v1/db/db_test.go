package db

import (
	"shop/pkg/options"
	"testing"
	"time"
)

func TestOrder_DB(t *testing.T) {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.16.192",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "shop_order_srv",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              1,
		EnableLog:             true,
	}
	_, err := NewOrderSQLClient(mysqlOpts)
	if err != nil {
		t.Error(err)
	}
}
