package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
)

func (gs *GoodsServer) CategoryBrandList(ctx context.Context, req *goods_pb.CategoryBrandFilterRequest) (*goods_pb.CategoryBrandListResponse, error) {
	var ret goods_pb.CategoryBrandListResponse

	list, err := gs.srv.CategoryBrand().List(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.CategoryBrandResponse{
			Id: item.ID,
			Brand: &goods_pb.BrandInfoResponse{
				Id:   item.Brands.ID,
				Name: item.Brands.Name,
				Logo: item.Brands.Logo,
			},
			Category: &goods_pb.CategoryInfoResponse{
				Id:             item.Category.ID,
				Name:           item.Category.Name,
				ParentCategory: item.Category.ParentCategoryID,
				Level:          item.Category.Level,
				IsTab:          item.Category.IsTab,
			},
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) GetCategoryBrandList(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.BrandListResponse, error) {
	var ret goods_pb.BrandListResponse

	list, err := gs.srv.CategoryBrand().Get(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.BrandInfoResponse{
			Id:   item.Brands.ID,
			Name: item.Brands.Name,
			Logo: item.Brands.Logo,
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) CreateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) (*goods_pb.CategoryBrandResponse, error) {
	var ret goods_pb.CategoryBrandResponse

	categoryBrandDO := do.CategoryBrandDO{
		CategoryID: req.CategoryId,
		BrandsID:   req.BrandId,
	}
	categoryBrandDTO := dto.CategoryBrandDTO{CategoryBrandDO: categoryBrandDO}

	categoryBrandID, err := gs.srv.CategoryBrand().Create(ctx, &categoryBrandDTO)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Id = categoryBrandID
	return &ret, nil
}

func (gs *GoodsServer) DeleteCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) (*empty.Empty, error) {
	err := gs.srv.CategoryBrand().Delete(ctx, req.Id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) UpdateCategoryBrand(ctx context.Context, req *goods_pb.CategoryBrandRequest) (*empty.Empty, error) {
	categoryBrandDO := do.CategoryBrandDO{
		CategoryID: req.CategoryId,
		BrandsID:   req.BrandId,
	}
	categoryBrandDTO := dto.CategoryBrandDTO{CategoryBrandDO: categoryBrandDO}

	if err := gs.srv.CategoryBrand().Update(ctx, &categoryBrandDTO); err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}
