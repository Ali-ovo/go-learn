package srvGoods

import (
	"context"
	goods_pb "shop/api/goods/v1"
)

type CategoryBrandSrv interface {
	// 品牌分类

	CategoryBrandList(ctx context.Context, req *goods_pb.CategoryBrandFilterRequest) (*goods_pb.CategoryBrandListResponse, error)

	// 通过category获取brands

	GetCategoryBrandList(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.BrandListResponse, error)
	CreateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) (*goods_pb.CategoryBrandResponse, error)
	DeleteCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) error
	UpdateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) error
}
