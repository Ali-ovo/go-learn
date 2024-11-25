package srvGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	dtoGoods "shop/app/shop/api/internal_api/domain/dto/goods"
)

type GoodsSrv interface {
	// 商品接口
	GoodsList(ctx context.Context, req *dtoGoods.GoodsFilter) (*dtoGoods.GoodDTOList, error)

	// 用户提交订单有多个商品，需要批量查询商品的信息

	BatchGetGoods(ctx context.Context, req *goods_pb.BatchGoodsIdInfo) (*goods_pb.GoodsListResponse, error)
	CreateGoods(ctx context.Context, req *goods_pb.CreateGoodsInfo) (*goods_pb.GoodsInfoResponse, error)
	DeleteGoods(ctx context.Context, req *goods_pb.DeleteGoodsInfo) error
	UpdateGoods(ctx context.Context, req *goods_pb.CreateGoodsInfo) error
	GetGoodsDetail(ctx context.Context, req *goods_pb.GoodInfoRequest) (*goods_pb.GoodsInfoResponse, error)
}
