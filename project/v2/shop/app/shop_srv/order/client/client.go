package main

import (
	"context"
	"fmt"
	"math/rand"
	order_pb "shop/api/order/v1"
	"shop/gmicro/registry/consul"
	rpc "shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"time"

	"github.com/hashicorp/consul/api"
)

// generateOrderSn 订单号的生成 (还有一种算法 雪花算法, 我们的订单号有一个问题, 不是递增的 如果需要递增 可以考虑雪花算法)
func generateOrderSn(userId int64) string {
	// 订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source) // 设置局部随机种子
	orderSn := fmt.Sprintf(
		"%04d%02d%02d%02d%02d%09d%09d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, randomGenerator.Intn(90)+10,
	)
	return orderSn
}

func main() {
	//设置全局的负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()

	conf := api.DefaultConfig()
	conf.Address = "172.16.89.133:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	conn, err := rpc.DialInsecure(
		context.Background(),
		rpc.WithBanlancerName("selector"),
		rpc.WithDiscovery(consul.New(cli, consul.WithHealthCheck(true))),
		rpc.WithEndpoint("discovery:///order_srv"),
		rpc.WithClientEnableTracing(false),
		rpc.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
		rpc.WithClientTimeout(time.Second*5000),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := order_pb.NewOrderClient(conn)

	_, err = uc.SubmitOrder(context.Background(), &order_pb.OrderRequest{
		UserId:  3,
		Address: "测试",
		Name:    "czc",
		Mobile:  "13067353692",
		Post:    "尽快发货",
		OrderSn: generateOrderSn(3),
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("订单创建成功!!")
}
