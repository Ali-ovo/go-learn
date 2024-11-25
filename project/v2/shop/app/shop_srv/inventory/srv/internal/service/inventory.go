package service

import (
	"context"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/app/shop_srv/inventory/srv/internal/domain/dto"
)

type InventorySrv interface {
	// Create 设置库存
	Create(ctx context.Context, inv *dto.InventoryDTO) error
	// Get 根据商品的id查询库存
	Get(ctx context.Context, goodsID int64) (*dto.InventoryDTO, error)
	// Sell 扣减库存
	Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error
	// Repack 归还库存
	Repack(ctx context.Context, ordersn string, details []do.GoodsDetail) error
}
