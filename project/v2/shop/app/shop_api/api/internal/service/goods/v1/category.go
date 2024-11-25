package serviceGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_api/api/internal/data/v1"
	srvGoods "shop/app/shop_api/api/internal/service/goods"
)

type CategoryService struct {
	data data.DataFactory
}

func (cs *CategoryService) GetAllCategorysList(ctx context.Context) (*goods_pb.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (cs *CategoryService) GetCategorysList(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (cs *CategoryService) GetSubCategory(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (cs *CategoryService) CreateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, req *goods_pb.DeleteCategoryRequest) error {
	//TODO implement me
	panic("implement me")
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewCategoryService(data data.DataFactory) srvGoods.CategorySrv {
	return &CategoryService{data: data}
}
