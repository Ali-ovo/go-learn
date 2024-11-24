package controller

import (
	"context"
	"fmt"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/data"
	"shop/app/shop_srv/goods/srv/internal/data/v1/db"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1/es"
	srv "shop/app/shop_srv/goods/srv/internal/service/v1"
	"shop/pkg/options"
	"testing"
)

func Conn() (data.DataFactory, data_search.SearchFactory) {
	// 数据库的工厂方法
	dataFactory, err := db.GetDBfactoryOr(&options.MySQLOptions{
		Host:      "192.168.101.49",
		Port:      3306,
		Username:  "root",
		Password:  "56248123",
		Database:  "shop_goods_srv",
		LogLevel:  4,
		EnableLog: true,
	})
	if err != nil {
		panic(err)
	}
	searchFactory, err := es.GetSearchFactoryOr(&options.EsOptions{
		Host:     "192.168.101.49",
		Port:     9200,
		Username: "elastic",
		Password: "56248Qwezxcv",
	})
	if err != nil {
		panic(err)
	}
	return dataFactory, searchFactory
}

func TestGoodsList(t *testing.T) {
	dataFactory, searchFactory := Conn()

	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, searchFactory)

	gServer := NewGoodsServer(srvFactory)

	list, err := gServer.GoodsList(context.Background(), &goods_pb.GoodsFilterRequest{
		PriceMin:    0,
		PriceMax:    0,
		IsHot:       false,
		IsNew:       false,
		IsTab:       false,
		TopCategory: 0,
		Pages:       0,
		PagePerNums: 0,
		KeyWords:    "猕猴桃",
		Brand:       0,
	})
	if err != nil {
		return
	}
	fmt.Println(list)
}

func TestGoodsCreate(t *testing.T) {
	dataFactory, searchFactory := Conn()

	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, searchFactory)

	gServer := NewGoodsServer(srvFactory)

	goods, err := gServer.CreateGoods(context.Background(), &goods_pb.CreateGoodsInfo{
		Name:            "测试2",
		GoodsSn:         "1654564564",
		Stocks:          1000,
		MarketPrice:     100,
		ShopPrice:       101,
		GoodsBrief:      "测试用2",
		ShipFree:        false,
		Images:          []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg"},
		DescImages:      []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg"},
		GoodsFrontImage: "http://xbfmxn.kw/uzhs",
		IsNew:           false,
		IsHot:           false,
		OnSale:          false,
		CategoryId:      130358,
		BrandId:         614,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(goods)
}

func TestGoodsUpdate(t *testing.T) {
	dataFactory, searchFactory := Conn()

	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, searchFactory)

	gServer := NewGoodsServer(srvFactory)

	goods, err := gServer.UpdateGoods(context.Background(), &goods_pb.CreateGoodsInfo{
		Id:              852,
		Name:            "测试",
		GoodsSn:         "1654564564",
		Stocks:          1000,
		MarketPrice:     100,
		ShopPrice:       101,
		GoodsBrief:      "测试用",
		ShipFree:        false,
		Images:          []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg"},
		DescImages:      []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg"},
		GoodsFrontImage: "http://xbfmxn.kw/uzhs",
		IsNew:           true,
		IsHot:           true,
		OnSale:          true,
		CategoryId:      130358,
		BrandId:         614,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(goods)
}

func TestGoodsDelete(t *testing.T) {
	dataFactory, searchFactory := Conn()

	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, searchFactory)

	gServer := NewGoodsServer(srvFactory)

	goods, err := gServer.DeleteGoods(context.Background(), &goods_pb.DeleteGoodsInfo{Id: 854})
	if err != nil {
		panic(err)
	}
	fmt.Println(goods)
}
