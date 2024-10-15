package global

import (
	"go-learn/shop/shop_api/order_web/config"
	"go-learn/shop/shop_api/order_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig        *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvClient     proto.GoodsClient
	OrderSrvClient     proto.OrderClient
	InventorySrvClient proto.InventoryClient
)
