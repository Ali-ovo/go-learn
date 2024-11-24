package es

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	searchFactory data_search.SearchFactory
	once          sync.Once
)

type dataSearch struct {
	esClient *elasticsearch.Client
}

func (ds *dataSearch) Goods() data_search.GoodsStore {
	return newGoods(ds)
}

func GetSearchFactoryOr(Esopts *options.EsOptions) (data_search.SearchFactory, error) {
	if Esopts == nil && searchFactory == nil {
		return nil, fmt.Errorf("es client is nil")
	}

	var err error
	once.Do(func() {
		esOpt := &conn.EsOptions{
			Host:     Esopts.Host,
			Port:     Esopts.Port,
			Username: Esopts.Username,
			Password: Esopts.Password,
		}
		esClient, err := conn.NewEsClient(esOpt)
		if err != nil {
			return
		}
		searchFactory = &dataSearch{esClient: esClient}
		log.Info("[goods-srv] 初始化 Es 完成")
	})

	if searchFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get es store factory")
	}
	return searchFactory, nil
}

type CanalData struct {
	Data      json.RawMessage  `json:"data"`
	Database  string           `json:"database"`
	Es        int64            `json:"es"`
	ID        int32            `json:"id"`
	IsDdl     bool             `json:"isDdl"`
	MysqlType map[string]any   `json:"mysqlType"`
	Old       []map[string]any `json:"old"`
	PkNames   []string         `json:"pkNames"`
	Sql       string           `json:"sql"`
	SqlType   map[string]any   `json:"sqlType"`
	Table     string           `json:"table"`
	Ts        uint64           `json:"ts"`
	Type      string           `json:"type"`
}

func GoodsSaveToES(Esopts *options.EsOptions) func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//searchFactory, err := GetSearchFactoryOr(Esopts)
		//if err != nil {
		//	log.Fatal(err.Error())
		//}

		for _, msg := range msgs {
			var data CanalData
			var goodsMap []map[string]interface{}
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				panic(err)
				return consumer.ConsumeRetryLater, err
			}

			if err := json.Unmarshal(data.Data, &goodsMap); err != nil {
				panic(err)
				return consumer.ConsumeRetryLater, err
			}

			switch data.Type {
			case "INSERT":
				//searchFactory.Goods().Update(ctx)
			case "DELETE":

			}

		}
		return consumer.ConsumeRetryLater, nil
	}
}
