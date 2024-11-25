package controller

import (
	"context"
	"fmt"
	inventory_pb "shop/api/inventory/v1"
	"shop/app/shop_srv/inventory/srv/internal/data/v1"
	"shop/app/shop_srv/inventory/srv/internal/data/v1/db"
	srv "shop/app/shop_srv/inventory/srv/internal/service/v1"
	"shop/gmicro/pkg/storage"
	"shop/pkg/options"
	"sync"
	"testing"
)

func Conn() data.DataFactory {
	dataFactory, err := db.GetDBfactoryOr(&options.MySQLOptions{
		Host:      "192.168.16.192",
		Port:      3306,
		Username:  "root",
		Password:  "56248123",
		Database:  "shop_inventory_srv",
		LogLevel:  4,
		EnableLog: true,
	})
	if err != nil {
		panic(err)
	}
	return dataFactory
}

func TestInventory_Sell(t *testing.T) {
	dataFactory := Conn()
	redisConfig := &storage.Config{
		Host:          "192.168.16.192",
		Port:          6379,
		EnableTracing: true,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go storage.ConnectToRedis(context.Background(), &wg, redisConfig)
	wg.Wait()

	srvFactory := srv.NewService(dataFactory)

	iServer := NewInventoryServer(srvFactory)

	goodsInv := []*inventory_pb.GoodsInvInfo{
		{
			GoodsId: 421,
			Num:     1,
		},
		{
			GoodsId: 422,
			Num:     2,
		},
		{
			GoodsId: 423,
			Num:     3,
		},
	}

	sell, err := iServer.Sell(context.Background(), &inventory_pb.SellInfo{
		GoodsInfo: goodsInv,
		OrderSn:   "222",
	})
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	fmt.Println(sell)
}

func TestInventory_Reback(t *testing.T) {
	dataFactory := Conn()
	redisConfig := &storage.Config{
		Host:          "192.168.16.192",
		Port:          6379,
		EnableTracing: true,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go storage.ConnectToRedis(context.Background(), &wg, redisConfig)
	wg.Wait()

	srvFactory := srv.NewService(dataFactory)

	iServer := NewInventoryServer(srvFactory)

	sell, err := iServer.Reback(context.Background(), &inventory_pb.SellInfo{
		//GoodsInfo: goodsInv,
		OrderSn: "222",
	})
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	fmt.Println(sell)
}
