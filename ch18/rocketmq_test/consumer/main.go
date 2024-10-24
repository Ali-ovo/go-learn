package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.189.128:9876"}),
		consumer.WithGroupName("ali1"),
	)

	err := c.Subscribe("ali1", consumer.MessageSelector{}, func(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range me {
			fmt.Printf("收到消息: %s\n", me[i])
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		panic("订阅消息失败")
	}

	_ = c.Start()

	time.Sleep(time.Hour)

	_ = c.Shutdown()

}
