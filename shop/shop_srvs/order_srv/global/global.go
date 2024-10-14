package global

import (
	"go-learn/shop/shop_srvs/order_srv/config"
	"go-learn/shop/shop_srvs/order_srv/proto"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
