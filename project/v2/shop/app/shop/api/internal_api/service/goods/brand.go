package srvGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
)

type BrandSrv interface {
	// 品牌和轮播图

	BrandList(ctx context.Context, req *goods_pb.BrandFilterRequest) (*goods_pb.BrandListResponse, error)
	CreateBrand(ctx context.Context, req *goods_pb.BrandRequest) (*goods_pb.BrandInfoResponse, error)
	DeleteBrand(ctx context.Context, req *goods_pb.BrandRequest) error
	UpdateBrand(ctx context.Context, req *goods_pb.BrandRequest) error
}
