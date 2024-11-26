package es

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/gmicro/pkg/mapstructure"
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
		searchFactory, err := GetSearchFactoryOr(Esopts)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, msg := range msgs {
			var data CanalData
			var mapGoods []map[string]interface{}
			var goodsSearchDOList []*do.GoodsSearchDO
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Error(err.Error())
				return consumer.ConsumeSuccess, err // 解析失败 说明 传递错误的数据进来
			}

			if err := json.Unmarshal(data.Data, &mapGoods); err != nil {
				log.Error(err.Error())
				return consumer.ConsumeRetryLater, err
			}

			err := mapstructure.Decode(mapGoods, &goodsSearchDOList)
			if err != nil {
				log.Error(err.Error())
				return consumer.ConsumeRetryLater, err
			}

			switch data.Type {
			case "INSERT":
				for _, goodsSearchDO := range goodsSearchDOList {
					//log.Infof("[goods-srv] 新增数据 %v", goodsSearchDO)
					err := searchFactory.Goods().Create(ctx, goodsSearchDO)
					if err != nil {
						log.Error(err.Error())
						return consumer.ConsumeRetryLater, err
					}
				}
			case "UPDATE":
				for _, goodsSearchDO := range goodsSearchDOList {
					if goodsSearchDO.DeleteAt == nil || !goodsSearchDO.DeleteAt.IsZero() {
						//log.Infof("[goods-srv] 删除数据 %v", goodsSearchDO)
						err := searchFactory.Goods().Delete(ctx, uint64(goodsSearchDO.ID))
						if err != nil {
							log.Error(err.Error())
							return consumer.ConsumeRetryLater, err
						}
					} else {
						//log.Infof("[goods-srv] 更新数据 %v", goodsSearchDO)
						err := searchFactory.Goods().Update(ctx, goodsSearchDO)
						if err != nil {
							log.Error(err.Error())
							return consumer.ConsumeRetryLater, err
						}
					}
				}
			case "DELETE":
				for _, goodsSearchDO := range goodsSearchDOList {
					//log.Infof("[goods-srv] 删除数据 %v", goodsSearchDO)
					err := searchFactory.Goods().Delete(ctx, uint64(goodsSearchDO.ID))
					if err != nil {
						log.Error(err.Error())
						return consumer.ConsumeRetryLater, err
					}
				}
			}
		}
		return consumer.ConsumeSuccess, nil
	}
}
