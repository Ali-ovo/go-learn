package service

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_srv/domain/dto"
)

type GoodsSrv interface {
	// List 商品列表
	List(ctx context.Context, req *goods_pb.GoodsFilterRequest) (*dto.GoodsDTOList, error)
	// Get 商品详情
	Get(ctx context.Context, id uint64) (*dto.GoodsDTO, error)
	// Create 创建商品
	Create(ctx context.Context, goods *dto.GoodsDTO) (int64, error)
	// Update 更新商品
	Update(ctx context.Context, goods *dto.GoodsDTO) error
	// Delete 删除商品
	Delete(ctx context.Context, id uint64) error
	// BatchGet 批量查询商品
	BatchGet(ctx context.Context, ids []int64) ([]*dto.GoodsDTO, error)
}
