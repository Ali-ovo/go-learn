package srv

import (
	"context"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data/v1/db"
	"shop/pkg/options"
	"testing"
	"time"
)

func conn() data.DataFactory {
	mysqlOpts := &options.MySQLOptions{
		Host:                  "192.168.189.128",
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
	dbFactory, err := db.GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}
	return dbFactory
}

func TestCategory_Get(t *testing.T) {
	dbFactory := conn()

	categoryService := CategoryService{
		data: dbFactory,
	}

	category, err := categoryService.Get(context.Background(), 130370)
	if err != nil {
		panic(err)
	}
	fmt.Println(category)
}

func TestCategory_List(t *testing.T) {
	dbFactory := conn()

	categoryService := CategoryService{
		data: dbFactory,
	}

	category, err := categoryService.List(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(category)
}
