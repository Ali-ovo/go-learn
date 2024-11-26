package data

import (
	goods_pb "shop/api/goods/v1"
	inventory_pb "shop/api/inventory/v1"

	"gorm.io/gorm"
)

type DataFactory interface {
	Orders() OrderStore
	ShopCarts() ShopCartStore
	Goods() goods_pb.GoodsClient
	Inventory() inventory_pb.InventoryClient
	Begin() *gorm.DB
}
