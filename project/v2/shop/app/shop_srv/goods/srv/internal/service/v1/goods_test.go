package srv

import (
	"context"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/data/v1/db"
	"shop/pkg/options"
	"testing"
	"time"
)

func TestGoodsBatchGet(t *testing.T) {
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
	dbFactory, err := db.GetDBfactoryOr(mysqlOpts)
	if err != nil {
		panic(err)
	}

	goodsService := goodsService{
		data: dbFactory,
	}

	get, err := goodsService.BatchGet(context.Background(), []uint64{840, 839, 838, 837, 836})
	if err != nil {
		panic(err)
	}

	//get, err = goodsService.BatchGetTwe(context.Background(), []uint64{845, 840, 839, 838, 837, 836})
	//if err != nil {
	//	panic(err)
	//}

	//get, err = goodsService.BatchGetThree(context.Background(), []uint64{845, 840, 839, 838, 837, 836})
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(get)
}
