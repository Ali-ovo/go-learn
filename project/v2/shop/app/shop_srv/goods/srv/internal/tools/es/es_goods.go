package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Mysql2Es() {
	dsn := "root:123456@tcp(192.168.189.128:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// 初始化连接
	host := "http://192.168.189.128:9200"
	esClient, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{host},
			Username:  "elastic",
			Password:  "56248Qwezxcv",
			//EnableDebugLogger: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// 新建 mapping 和 index
	searchDO := do.GoodsSearchDO{}
	exists, err := esapi.IndicesExistsRequest{Index: []string{searchDO.GetIndexName()}}.Do(context.Background(), esClient)
	if err != nil {
		panic(err)
	}
	defer exists.Body.Close()
	if exists.IsError() {
		// Create a new index.
		exists, err := esapi.IndicesCreateRequest{
			Index: searchDO.GetIndexName(),
			Body:  bytes.NewReader([]byte(searchDO.GetMapping())),
		}.Do(context.Background(), esClient)
		if err != nil {
			// Handle error
			panic(err)
		}
		defer exists.Body.Close()
	}

	var goods []do.GoodsDO
	db.Find(&goods)
	for _, g := range goods {
		esModel := do.GoodsSearchDO{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandsID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}

		res, err := esapi.IndexRequest{
			Index:      esModel.GetIndexName(),
			DocumentID: strconv.Itoa(int(g.ID)),
			Body:       esutil.NewJSONReader(esModel),
		}.Do(context.Background(), esClient)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		//强调一下 一定要将docker启动es的java_ops的内存设置大一些 否则运行过程中会出现 bad request错误
	}
}

func main() {
	Mysql2Es()
}
