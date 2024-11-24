package db

import (
	"context"
	"fmt"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestBrands_Get(t *testing.T) {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.101.49",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "shop_goods_srv",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              1,
		EnableLog:             true,
	}
	dbFactory, err := GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}

	rsp, err := dbFactory.Brands().Get(context.Background(), 111111)
	if err != nil || rsp == nil {
		panic(err)
	}

	fmt.Println(rsp)
}