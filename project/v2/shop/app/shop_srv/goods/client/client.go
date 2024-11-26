package main

import (
	"context"
	"fmt"
	goods_pb "shop/api/goods/v1"
	"shop/gmicro/registry/consul"
	rpc "shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"time"

	"github.com/hashicorp/consul/api"
)

func main() {
	//设置全局的负载均衡策略
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()

	conf := api.DefaultConfig()
	conf.Address = "192.168.189.128:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	conn, err := rpc.DialInsecure(
		context.Background(),
		rpc.WithBanlancerName("selector"),
		rpc.WithDiscovery(consul.New(cli, consul.WithHealthCheck(true))),
		rpc.WithEndpoint("discovery:///goods_srv"),
		rpc.WithClientEnableTracing(false),
		rpc.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
		rpc.WithClientTimeout(time.Second*5000),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := goods_pb.NewGoodsClient(conn)

	re, err := uc.GoodsList(context.Background(), &goods_pb.GoodsFilterRequest{
		KeyWords: "猕猴桃",
	})
	if err != nil {
		panic(err)
	}

	for _, item := range re.Data {
		fmt.Println(item)
	}
	time.Sleep(time.Second * 3)
}
