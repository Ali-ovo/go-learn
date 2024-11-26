package db

import (
	"context"
	"fmt"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestBrands_Get(t *testing.T) {
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

	rsp, err := dbFactory.Brands().Get(context.Background(), 111111)
	if err != nil || rsp == nil {
		panic(err)
	}

	fmt.Println(rsp)
}

func TestBrands_List(t *testing.T) {
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
		t.Error(err.Error())
	}

	rsp, err := dbFactory.Brands().List(context.Background(), metav1.ListMeta{}, []string{})
	if err != nil || rsp == nil {
		t.Error(err.Error())
	}

	fmt.Println(rsp)
}
