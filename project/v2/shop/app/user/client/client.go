package main

import (
	"context"
	"fmt"
	v1 "shop/api/user/v1"
	rpc "shop/gmicro/server/rpcserver"
	_ "shop/gmicro/server/rpcserver/resolver/direct" // 必填

	"github.com/hashicorp/consul/api"

	"shop/gmicro/registry/consul"
)

func main() {

	conf := api.DefaultConfig()
	conf.Address = "172.16.89.133:8500"
	conf.Scheme = "http"
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	conn, err := rpc.DialInsecure(context.Background(),
		// rpc.WithEndpoint("127.0.0.1:8081"),
		// rpc.WithEndpoint("discovery:///172.16.100.208:8081"),
		// rpc.WithClientTimeout(1*time.Second),

		rpc.WithDiscovery(consul.New(cli)),
		rpc.WithEndpoint("discovery:///grpcServer"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	uc := v1.NewUserClient(conn)
	re, err := uc.GetUserList(context.Background(), &v1.PageInfo{})
	if err != nil {
		panic(err)
	}
	fmt.Println(re)
}
