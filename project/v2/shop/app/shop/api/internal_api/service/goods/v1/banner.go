package serviceGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_api/data/v1"
	srvGoods "shop/app/shop/api/internal_api/service/goods"
)

type BannerService struct {
	data data.DataFactory
}

func (bs *BannerService) BannerList(ctx context.Context) (*goods_pb.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (bs *BannerService) CreateBanner(ctx context.Context, req *goods_pb.BannerRequest) (*goods_pb.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (bs *BannerService) DeleteBanner(ctx context.Context, req *goods_pb.BannerRequest) error {
	//TODO implement me
	panic("implement me")
}

func (bs *BannerService) UpdateBanner(ctx context.Context, req *goods_pb.BannerRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewBannerService(data data.DataFactory) srvGoods.BannerSrv {
	return &BannerService{data: data}
}
