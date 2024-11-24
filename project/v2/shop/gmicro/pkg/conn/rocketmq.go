package conn

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

// ConsumerOptions PS: 具体参数有需要再看 源码 添加
type ConsumerOptions struct {
	Addr      []string `mapstructure:"addr" json:"addr,omitempty"`
	GroupName string   `mapstructure:"group_name" json:"group_name,omitempty"`
}

func NewRocketmqConsumer(opts *ConsumerOptions) (rocketmq.PushConsumer, error) {
	rlog.SetLogLevel("error")
	return rocketmq.NewPushConsumer(
		consumer.WithGroupName(opts.GroupName),
		consumer.WithNsResolver(
			primitive.NewPassthroughResolver(opts.Addr),
		),
	)
}
