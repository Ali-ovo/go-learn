package srv

import (
	"shop/app/shop_srv/goods/srv/config"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1/es"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func RocketmqConsumer(opts *config.Config) {
	var err error

	rocketMQ := conn.RocketMQ{
		Addr:          opts.Rocketmq.Addr,
		ConsumerGroup: opts.Rocketmq.GroupName,
		Retry:         opts.Rocketmq.Retry,
		Consumer: []conn.Consumer{
			{
				Topic:    "goods_canal",
				Selector: consumer.MessageSelector{},
				Func:     es.GoodsSaveToES(opts.EsOptions),
			},
		},
	}

	// 连接 rockerMq 启动消费者消费消息
	err = rocketMQ.ConsumerStart()
	if err != nil {
		panic(err)
	}
	log.Info("[goods-srv] 启动 Rocketmq 完成")
}
