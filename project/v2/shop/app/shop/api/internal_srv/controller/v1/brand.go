package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_srv/domain/do"
	"shop/app/shop/api/internal_srv/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/pkg/gorm"

	"github.com/golang/protobuf/ptypes/empty"
)

func (gs *GoodsServer) BrandList(ctx context.Context, req *goods_pb.BrandFilterRequest) (*goods_pb.BrandListResponse, error) {
	var ret goods_pb.BrandListResponse

	list, err := gs.srv.Brand().List(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.BrandInfoResponse{
			Id:   item.ID,
			Name: item.Name,
			Logo: item.Logo,
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) CreateBrand(ctx context.Context, req *goods_pb.BrandRequest) (*goods_pb.BrandInfoResponse, error) {
	var ret goods_pb.BrandInfoResponse

	brandsDO := do.BrandsDO{
		Name: req.Name,
		Logo: req.Logo,
	}
	brandsDTO := dto.BrandsDTO{brandsDO}

	brandID, err := gs.srv.Brand().Create(ctx, &brandsDTO)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Id = brandID
	return &ret, nil
}

func (gs *GoodsServer) DeleteBrand(ctx context.Context, req *goods_pb.BrandRequest) (*empty.Empty, error) {
	err := gs.srv.Brand().Delete(ctx, req.Id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) UpdateBrand(ctx context.Context, req *goods_pb.BrandRequest) (*empty.Empty, error) {
	brandDO := do.BrandsDO{
		BaseModel: gorm.BaseModel{ID: req.Id},
		Name:      req.Name,
		Logo:      req.Logo,
	}
	brandDTO := dto.BrandsDTO{BrandsDO: brandDO}

	if err := gs.srv.Brand().Update(ctx, &brandDTO); err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}
