package srvGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
)

type CategorySrv interface {
	// 商品分类

	GetAllCategorysList(ctx context.Context) (*goods_pb.CategoryListResponse, error)
	GetCategorysList(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.CategoryListResponse, error)
	GetSubCategory(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.SubCategoryListResponse, error)
	CreateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.CategoryInfoResponse, error)
	DeleteCategory(ctx context.Context, req *goods_pb.DeleteCategoryRequest) error
	UpdateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) error
}
