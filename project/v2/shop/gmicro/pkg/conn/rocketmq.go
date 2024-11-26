package conn

import (
	"context"
	"shop/gmicro/pkg/errors"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

var (
	ConsumerOnce   sync.Once
	ConsumerClient rocketmq.PushConsumer
	ProducerOnce   sync.Once
	ProducerClient rocketmq.Producer
)

type RocketMQProducer interface {
	ProducerStart() error // 开启生产者Client
	Message() error       // 发送消息
}

type RocketMQConsumer interface {
	Subscribe() error     // 订阅消息
	ConsumerStart() error // 开启消费者Client
}

type Consumer struct {
	Topic    string
	Selector consumer.MessageSelector                                                        // 对 消息的属性进行条件过滤，例如标签、键值对 等
	Func     func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error) // 可能存在多个订阅  使用
}

type Producer struct {
	Topic          string
	Data           []byte
	DelayTimeLevel int
}

func init() {
	rlog.SetLogLevel("error")
}

// RocketMQ Options PS: 具体参数有需要再看 源码 添加
// TODO Rocketmq 事务消息未实现 因为我现在用不到
type RocketMQ struct {
	Addr          []string `mapstructure:"addr" json:"addr,omitempty"`
	ConsumerGroup string   `mapstructure:"consumer_group" json:"consumer_group,omitempty"`
	ProducerGroup string   `mapstructure:"producer_group" json:"producer_group"`
	Retry         int      `mapstructure:"retry" json:"retry,omitempty"` // 发送到 rocketmq 重试次数
	Consumer      []Consumer
	Producer      Producer
}

func (rmq *RocketMQ) ConsumerStart() error {
	var err error

	ConsumerOnce.Do(func() {
		var consumerOpts []consumer.Option

		if len(rmq.Addr) > 0 {
			consumerOpts = append(consumerOpts, consumer.WithNsResolver(
				primitive.NewPassthroughResolver(rmq.Addr),
			))
		} else {
			err = errors.Errorf("rocketmq Addr not found")
			return
		}

		if rmq.ConsumerGroup != "" {
			consumerOpts = append(consumerOpts, consumer.WithGroupName(rmq.ConsumerGroup))
		}

		if rmq.Retry > 0 {
			consumerOpts = append(consumerOpts, consumer.WithRetry(rmq.Retry))
		}

		ConsumerClient, err = rocketmq.NewPushConsumer(consumerOpts...)
		if err != nil {
			return
		}

		err = rmq.Subscribe()
		if err != nil {
			return
		}

		if err = ConsumerClient.Start(); err != nil {
			return
		}
	})
	if ProducerClient == nil || err != nil {
		return err
	}
	return nil
}

func (rmq *RocketMQ) Subscribe() error {
	var err error

	for _, subscribe := range rmq.Consumer {
		err = ConsumerClient.Subscribe(
			subscribe.Topic,
			subscribe.Selector,
			subscribe.Func,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rmq *RocketMQ) ProducerStart() error {
	var err error

	ProducerOnce.Do(func() {
		var producerOpts []producer.Option

		if len(rmq.Addr) > 0 {
			producerOpts = append(producerOpts, producer.WithNsResolver(
				primitive.NewPassthroughResolver(rmq.Addr),
			))
		} else {
			err = errors.Errorf("rocketmq Addr not found")
			return
		}

		if rmq.ConsumerGroup != "" {
			producerOpts = append(producerOpts, producer.WithGroupName(rmq.ConsumerGroup))
		}

		if rmq.Retry > 0 {
			producerOpts = append(producerOpts, producer.WithRetry(rmq.Retry))
		}

		ProducerClient, err = rocketmq.NewProducer(producerOpts...)
		if err != nil {
			return
		}
		if err = ProducerClient.Start(); err != nil {
			return
		}
	})
	if ProducerClient == nil || err != nil {
		return err
	}
	return nil
}

func (rmq *RocketMQ) Message() error {
	var err error
	// 发送 消息
	msg := primitive.NewMessage(rmq.Producer.Topic, rmq.Producer.Data)
	if rmq.Producer.DelayTimeLevel != 0 {
		msg.WithDelayTimeLevel(rmq.Producer.DelayTimeLevel)
	}
	_, err = ProducerClient.SendSync(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}

func Message(topic string, data []byte, delayTimeLevel int) error {
	var err error
	// 发送 消息
	msg := primitive.NewMessage(topic, data)
	if delayTimeLevel != 0 {
		msg.WithDelayTimeLevel(delayTimeLevel)
	}
	_, err = ProducerClient.SendSync(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}
