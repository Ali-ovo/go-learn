package db

import (
	"fmt"
	goods_pb "shop/api/goods/v1"
	inventory_pb "shop/api/inventory/v1"
	"shop/app/shop_srv/order/srv/config"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/gmicro/registry"
	"shop/gmicro/registry/consul"
	"shop/gmicro/server/rpcserver"
	"shop/gmicro/server/rpcserver/selector"
	"shop/gmicro/server/rpcserver/selector/random"
	"shop/pkg/options"
	"sync"

	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var once sync.Once

type DataFactory struct {
	db              *gorm.DB
	goodsClient     goods_pb.GoodsClient
	inventoryClient inventory_pb.InventoryClient
}

func (df *DataFactory) Orders() data.OrderStore {
	panic("not implemented")
}

func (df *DataFactory) ShopCarts() data.ShopCartStore {
	panic("not implemented")
}

func (df *DataFactory) Goods() goods_pb.GoodsClient {
	//TODO implement me
	panic("implement me")
}

func (df *DataFactory) Inventory() inventory_pb.InventoryClient {
	//TODO implement me
	panic("implement me")
}

func (df *DataFactory) Begin() *gorm.DB {
	return df.db.Begin()
}

func GetDataFactoryOr(options *config.Config) (data.DataFactory, error) {
	if options == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get grpc store factory")
	}

	var err error

	once.Do(func() {
		var gormDB *gorm.DB
		var InventoryClient inventory_pb.InventoryClient
		var goodsClient goods_pb.GoodsClient

		gormDB, err = NewOrderSQLClient(options.Mysql)
		if err != nil {
			return
		}

		selector.SetGlobalSelector(random.NewBuilder()) // 设置全局的负载均衡策略
		rpcserver.InitBuilder()                         // 构建负载均衡器的 构建器

		discovery := NewDiscovery(options.Registry)
		InventoryClient, err = NewInventoryServiceClient(discovery)
		goodsClient, err = NewGoodsServiceClient(discovery)
		dbFactory = &DataFactory{gormDB, goodsClient, InventoryClient}
	})

	if dbFactory == nil || err != nil {
		return nil, err
	}
	return dbFactory, nil
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
