package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct {
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("执行本地事务")
	time.Sleep(time.Second * 3)
	fmt.Println("本地事务执行失败")
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("检查本地事务")
	time.Sleep(time.Second * 15)
	return primitive.CommitMessageState

}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"172.16.89.133:9876"}))

	if err != nil {
		panic("生成 producer 失败")
	}

	if err = p.Start(); err != nil {
		panic("启动 producer 失败")
	}

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TransTopic", []byte("Hello transaction RocketMQ222")))

	if err != nil {
		fmt.Printf("发送消息失败: %s\n", err)
	} else {
		fmt.Printf("发送消息成功: %s\n", res.String())
	}

	time.Sleep(time.Hour)

	if err = p.Shutdown(); err != nil {
		panic("关闭 producer 失败")
	}
}
