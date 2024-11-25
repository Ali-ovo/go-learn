package data

import (
	"context"
	"database/sql"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
)

type InventoryStore interface {
	// Get 查询商品的库存信息
	Get(ctx context.Context, goodsID int64) (*do.InventoryDO, error)
	// Create 新建库存信息
	Create(ctx context.Context, inventoryDO *do.InventoryDO) error
	// GetSellDetail 查询库存销售信息
	GetSellDetail(ctx context.Context, ordersn string) (*do.StockSellDetailDO, error)
	// Reduce 扣减库存
	Reduce(ctx context.Context, txn *sql.Tx, goodsID int64, num int32) (sql.Result, error)
	// Increase 新增库存
	Increase(ctx context.Context, txn *sql.Tx, goodsID int64, num int32) (sql.Result, error)
	// CreateStockSellDetail 创建 订单
	CreateStockSellDetail(ctx context.Context, txn *sql.Tx, detail *do.StockSellDetailDO) (sql.Result, error)
	// UpdateStockSellDetailStatus 修改 订单状态
	UpdateStockSellDetailStatus(ctx context.Context, txn *sql.Tx, ordersn string, status int32) (sql.Result, error)
}
