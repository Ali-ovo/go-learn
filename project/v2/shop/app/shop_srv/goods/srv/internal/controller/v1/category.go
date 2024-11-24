package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
)

func (gs *GoodsServer) GetAllCategorysList(ctx context.Context, e *empty.Empty) (*goods_pb.CategoryListResponse, error) {
	var ret goods_pb.CategoryListResponse

	list, err := gs.srv.Category().AllList(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.CategoryInfoResponse{
			Id:             item.ID,
			Name:           item.Name,
			ParentCategory: item.ParentCategoryID,
			Level:          item.Level,
			IsTab:          item.IsTab,
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) GetCategorysList(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.CategoryListResponse, error) {
	var ret goods_pb.CategoryListResponse

	list, err := gs.srv.Category().List(ctx, req.Level)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.CategoryInfoResponse{
			Id:             item.ID,
			Name:           item.Name,
			ParentCategory: item.ParentCategoryID,
			Level:          item.Level,
			IsTab:          item.IsTab,
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) GetSubCategory(ctx context.Context, req *goods_pb.CategoryListRequest) (*goods_pb.SubCategoryListResponse, error) {
	var ret goods_pb.SubCategoryListResponse

	category, err := gs.srv.Category().Get(ctx, int64(req.Id))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	ret.Info = &goods_pb.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}
	for _, item := range category.SubCategory {
		ret.SubCategorys = append(ret.SubCategorys, &goods_pb.CategoryInfoResponse{
			Id:             item.ID,
			Name:           item.Name,
			ParentCategory: item.ParentCategoryID,
			Level:          item.Level,
			IsTab:          item.IsTab,
		})
	}
	ret.Total = int32(len(ret.SubCategorys))
	return &ret, nil
}

func (gs *GoodsServer) CreateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*goods_pb.CategoryInfoResponse, error) {
	var ret goods_pb.CategoryInfoResponse

	categoryDO := do.CategoryDO{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategory,
		Level:            req.Level,
		IsTab:            req.IsTab,
	}
	categoryDTO := dto.CategoryDTO{CategoryDO: categoryDO}

	categoryID, err := gs.srv.Category().Create(ctx, &categoryDTO)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Id = categoryID
	return &ret, nil
}

func (gs *GoodsServer) DeleteCategory(ctx context.Context, req *goods_pb.DeleteCategoryRequest) (*empty.Empty, error) {
	err := gs.srv.Category().Delete(ctx, int64(req.Id))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) UpdateCategory(ctx context.Context, req *goods_pb.CategoryInfoRequest) (*empty.Empty, error) {
	categoryDO := do.CategoryDO{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategory,
		Level:            req.Level,
		IsTab:            req.IsTab,
	}
	categoryDTO := dto.CategoryDTO{CategoryDO: categoryDO}

	if err := gs.srv.Category().Update(ctx, &categoryDTO); err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}
