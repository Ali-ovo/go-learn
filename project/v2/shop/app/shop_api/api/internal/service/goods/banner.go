package srvGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
)

type BannerSrv interface {
	// 轮播图

	BannerList(ctx context.Context) (*goods_pb.BannerListResponse, error)
	CreateBanner(ctx context.Context, req *goods_pb.BannerRequest) (*goods_pb.BannerResponse, error)
	DeleteBanner(ctx context.Context, req *goods_pb.BannerRequest) error
	UpdateBanner(ctx context.Context, req *goods_pb.BannerRequest) error
}
