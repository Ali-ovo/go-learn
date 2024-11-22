package main

import (
	"context"
	"fmt"
	v1 "shop/api/user/v1"
	rpc "shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	_ "shop/gmicro/server/rpcserver/resolver/direct" // 必填
	"time"

	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"

	"github.com/hashicorp/consul/api"

	"shop/gmicro/registry/consul"
)

func main() {

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
		rpc.WithDiscovery(consul.New(cli, consul.WithHealthCheck(true))),
		rpc.WithEndpoint("discovery:///api_user_srv"),
		rpc.WithClientEnableTracing(false),
		rpc.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
		rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)

	for {
		r, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
		if err != nil {
			panic(err)
		}
		fmt.Println(r)
		time.Sleep(time.Second * 5)
	}
}