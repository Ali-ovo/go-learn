package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/pkg/gorm"

	"github.com/golang/protobuf/ptypes/empty"
)

func (gs *GoodsServer) BannerList(ctx context.Context, e *empty.Empty) (*goods_pb.BannerListResponse, error) {
	var ret goods_pb.BannerListResponse

	list, err := gs.srv.Banner().List(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, &goods_pb.BannerResponse{
			Id:    item.ID,
			Image: item.Image,
			Index: item.Index,
			Url:   item.Url,
		})
	}
	return &ret, nil
}

func (gs *GoodsServer) CreateBanner(ctx context.Context, request *goods_pb.BannerRequest) (*goods_pb.BannerResponse, error) {
	var ret goods_pb.BannerResponse

	bannerDO := do.BannerDO{
		Image: request.Image,
		Url:   request.Url,
		Index: request.Index,
	}
	bannerDTO := dto.BannerDTO{BannerDO: bannerDO}

	bannerID, err := gs.srv.Banner().Create(ctx, &bannerDTO)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Id = bannerID
	return &ret, nil
}

func (gs *GoodsServer) DeleteBanner(ctx context.Context, req *goods_pb.BannerRequest) (*empty.Empty, error) {
	err := gs.srv.Banner().Delete(ctx, req.Id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) UpdateBanner(ctx context.Context, req *goods_pb.BannerRequest) (*empty.Empty, error) {
	bannerDO := do.BannerDO{
		BaseModel: gorm.BaseModel{ID: req.Id},
		Image:     req.Image,
		Url:       req.Url,
		Index:     req.Index,
	}
	bannerDTO := dto.BannerDTO{BannerDO: bannerDO}

	if err := gs.srv.Banner().Update(ctx, &bannerDTO); err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}
