package serviceGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_api/data/v1"
	srvGoods "shop/app/shop/api/internal_api/service/goods"
)

type CategoryBrandService struct {
	data data.DataFactory
}

func (c CategoryBrandService) CategoryBrandList(ctx context.Context, req *goods_pb.CategoryBrandFilterRequest) (*goods_pb.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryBrandService) GetCategoryBrandList(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryBrandService) CreateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) (*goods_pb.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryBrandService) DeleteCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c CategoryBrandService) UpdateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewCategoryBrandService(data data.DataFactory) srvGoods.CategoryBrandSrv {
	return &CategoryBrandService{data: data}
}
