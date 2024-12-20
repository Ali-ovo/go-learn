package main

import (
	"context"
	"fmt"
	user_pb "shop/api/user/v1"
	"shop/gmicro/registry/consul"
	rpc "shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/clientinterceptors"

	"github.com/hashicorp/consul/api"

	//_ "shop/gmicro/server/rpcserver/resolver/direct" // 必填
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"time"
)

func main() {
	// 设置全局的负载均衡策略
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
		// 设置负载均衡
		rpc.WithBanlancerName("selector"),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpc.WithDiscovery(consul.New(cli)),
		rpc.WithEndpoint("discovery:///user_srv"),
		rpc.WithClientEnableTracing(false),
		rpc.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
		rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := user_pb.NewUserClient(conn)
	for {
		_, err := uc.GetUserList(context.Background(), &user_pb.PageInfo{})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("success")
		time.Sleep(time.Millisecond * 2)
	}

}
