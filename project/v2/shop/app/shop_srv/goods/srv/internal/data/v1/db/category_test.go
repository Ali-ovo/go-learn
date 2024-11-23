package db

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestListAll(t *testing.T) {
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

	category := NewCategory(db)
	rsp, err := category.Get(context.Background(), 130364)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(&rsp)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))

}
