package serviceGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_api/data/v1"
	srvGoods "shop/app/shop/api/internal_api/service/goods"
)

type BrandService struct {
	data data.DataFactory
}

func (bs *BrandService) BrandList(ctx context.Context, req *goods_pb.BrandFilterRequest) (*goods_pb.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (bs *BrandService) CreateBrand(ctx context.Context, req *goods_pb.BrandRequest) (*goods_pb.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (bs *BrandService) DeleteBrand(ctx context.Context, req *goods_pb.BrandRequest) error {
	//TODO implement me
	panic("implement me")
}

func (bs *BrandService) UpdateBrand(ctx context.Context, req *goods_pb.BrandRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewBrandService(data data.DataFactory) srvGoods.BrandSrv {
	return &BrandService{data: data}
}
