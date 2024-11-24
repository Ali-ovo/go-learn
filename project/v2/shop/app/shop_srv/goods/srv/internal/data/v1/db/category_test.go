package db

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/data"
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
	return dbFactory
}

func TestCategory_ListAll(t *testing.T) {
	dbFactory := Conn()
	rsp, err := dbFactory.Category().Get(context.Background(), 130364)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(&rsp)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

func TestCategory_List(t *testing.T) {
	dbFactory := Conn()
	rsp, err := dbFactory.Category().List(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(&rsp)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
