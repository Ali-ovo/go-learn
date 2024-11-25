package srv

import (
	"shop/app/shop_srv/inventory/srv/internal/data/v2"
	"shop/app/shop_srv/inventory/srv/internal/service"
	"shop/gmicro/pkg/storage"

	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

type serviceFactory struct {
	data data.DataFactory
	pool redsyncredis.Pool
}

func (sf *serviceFactory) Inventory() service.InventorySrv {
	return newInventory(sf)
}

func NewService(store data.DataFactory) service.ServiceFactory {
	rstore := storage.RedisCluster{}
	pool := goredis.NewPool(rstore.GetClient())
	return &serviceFactory{data: store, pool: pool}
}
