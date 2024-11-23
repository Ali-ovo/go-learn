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
		Host:                  "192.168.189.128",
		Port:                  3306,
		Username:              "root",
		Password:              "56248123",
		Database:              "mxshop_goods_srv",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Second * time.Duration(10),
		LogLevel:              1,
		EnableLog:             true,
	}
	db, err := GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}

	brands := NewBrand(db)
	rsp, err := brands.Get(context.Background(), 111111)
	if err != nil || rsp == nil {
		panic(err)
	}

	fmt.Println(rsp)
}
