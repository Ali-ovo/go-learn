package db

import (
	"context"
	"fmt"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestList(t *testing.T) {
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
	dbFactory, err := GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}

	rsp, err := dbFactory.Brands().List(context.Background(), metav1.ListMeta{
		Page:     0,
		PageSize: 10,
	}, []string{})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.TotalCount, rsp.Items)
}
