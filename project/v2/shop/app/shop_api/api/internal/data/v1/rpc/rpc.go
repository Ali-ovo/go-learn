package rpc

import (
	"shop/app/shop_api/api/internal/data/v1"
	"shop/gmicro/registry"
	"shop/gmicro/registry/consul"
	"shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"shop/pkg/options"

	"github.com/hashicorp/consul/api"
)

func NewDiscovery(options *options.RegistryOptions) registry.Discovery {
	conf := api.DefaultConfig()
	conf.Address = options.Address
	conf.Scheme = options.Scheme
	cli, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}
	return consul.New(cli, consul.WithHealthCheck(true))
}

// GetDataFactoryOr gRPC 的连接  基于服务发现
func GetDataFactoryOr(options *options.RegistryOptions) (data.UserData, error) {
	// 这里负责底层依赖的所有的rpc连接
	selector.SetGlobalSelector(random.NewBuilder()) // 设置全局的负载均衡策略
	rpcserver.InitBuilder()

	discovery := NewDiscovery(options)
	userClient := NewUserServiceClient(discovery)
	return NewUsers(userClient), nil
}
