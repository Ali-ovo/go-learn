package data

import (
	"context"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"

	"gorm.io/gorm"
)

type InventoryStore interface {
	// Get 查询商品的库存信息
	Get(ctx context.Context, goodsID int64) (*do.InventoryDO, error)
	// Create 新建库存信息
	Create(ctx context.Context, txn *gorm.DB, inventoryDO *do.InventoryDO) *gorm.DB
	// GetSellDetail 查询库存销售信息
	GetSellDetail(ctx context.Context, ordersn string) (*do.StockSellDetailDO, error)
	// Reduce 扣减库存
	Reduce(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) *gorm.DB
	// Increase 新增库存
	Increase(ctx context.Context, txn *gorm.DB, goodsID int64, num int32) *gorm.DB
	// CreateStockSellDetail 创建 订单
	CreateStockSellDetail(ctx context.Context, txn *gorm.DB, detail *do.StockSellDetailDO) *gorm.DB
	// UpdateStockSellDetailStatus 修改 订单状态
	UpdateStockSellDetailStatus(ctx context.Context, txn *gorm.DB, ordersn string, status int32) *gorm.DB
}
