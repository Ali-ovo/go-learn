package rpc

import (
	"fmt"
	goods_pb "shop/api/goods/v1"
	user_pb "shop/api/user/v1"
	"shop/app/shop_api/api/internal/data/v2"
	"shop/app/shop_api/api/internal/data/v2/goodsSrv"
	"shop/app/shop_api/api/internal/data/v2/rpc/goods"
	"shop/app/shop_api/api/internal/data/v2/rpc/user"
	"shop/app/shop_api/api/internal/data/v2/userSrv"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/registry"
	"shop/gmicro/registry/consul"
	"shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"shop/pkg/code"
	"shop/pkg/options"
	"sync"

	"github.com/hashicorp/consul/api"
)

var (
	dbFactory data.DataFactory
	once      sync.Once
)

type GrpcFactory struct {
	UserClient  user_pb.UserClient
	GoodsClient goods_pb.GoodsClient
}

func (g GrpcFactory) User() userSrv.UserData {
	return user.NewUser(g)
}

func (g GrpcFactory) Goods() goodsSrv.GoodsData {
	return goods.NewGoods(g)
}

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
func GetDataFactoryOr(options *options.RegistryOptions) (data.DataFactory, error) {
	if options == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get grpc store factory")
	}

	// 这里负责底层依赖的所有的rpc连接
	once.Do(func() {
		selector.SetGlobalSelector(random.NewBuilder()) // 设置全局的负载均衡策略
		rpcserver.InitBuilder()                         // 构建负载均衡器的 构建器

		discovery := NewDiscovery(options)
		userClient := user.NewUserServiceClient(discovery)
		goodsClient := goods.NewGoodsServiceClient(discovery)
		dbFactory = &GrpcFactory{userClient, goodsClient}
	})

	if dbFactory == nil {
		return nil, errors.WithCode(code.ErrConnectGRPC, "failed to get grpc store factory")
	}
	return dbFactory, nil
}
