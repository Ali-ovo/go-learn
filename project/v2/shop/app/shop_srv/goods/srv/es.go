package srv

import (
	"fmt"
	"shop/app/shop_srv/goods/srv/config"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1/es"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func RocketmqToEs(opts *config.Config) {
	client, err := conn.NewRocketmqConsumer(&conn.ConsumerOptions{
		Addr:      opts.Rocketmq.Addr,
		GroupName: opts.Rocketmq.GroupName,
	})
	if err != nil {
		panic(err)
	}

	client.Subscribe(
		"goods_canal",
		consumer.MessageSelector{},
		es.GoodsSaveToES(opts.EsOptions),
	)
	if err = client.Start(); err != nil {
		panic(fmt.Sprintf("启动 producer 失败: %s", err))
	}
	log.Info("[goods-srv] 启动 Rocketmq 完成")
}
