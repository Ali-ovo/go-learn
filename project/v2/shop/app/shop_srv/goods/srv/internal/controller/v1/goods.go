package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
)

func (gs *GoodsServer) GoodsList(ctx context.Context, request *goods_pb.GoodsFilterRequest) (*goods_pb.GoodsListResponse, error) {
	var ret goods_pb.GoodsListResponse

	list, err := gs.srv.Goods().List(ctx, request)
	if err != nil {
		//log.Errorf("get goods list error: %v", err.Error())
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(list.TotalCount)
	for _, item := range list.Items {
		ret.Data = append(ret.Data, DTOToResponse(item))
	}
	return &ret, nil
}

func (gs *GoodsServer) BatchGetGoods(ctx context.Context, info *goods_pb.BatchGoodsIdInfo) (*goods_pb.GoodsListResponse, error) {
	var ret goods_pb.GoodsListResponse

	list, err := gs.srv.Goods().BatchGet(ctx, info.Id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	ret.Total = int32(len(list))
	for _, item := range list {
		ret.Data = append(ret.Data, DTOToResponse(item))
	}
	return &ret, nil
}

func (gs *GoodsServer) CreateGoods(ctx context.Context, info *goods_pb.CreateGoodsInfo) (*goods_pb.GoodsInfoResponse, error) {
	var ret goods_pb.GoodsInfoResponse

	goodsDO := do.GoodsDO{
		CategoryID:      info.CategoryId,
		BrandsID:        info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}
	goodsDTO := dto.GoodsDTO{GoodsDO: goodsDO}

	goodsID, err := gs.srv.Goods().Create(ctx, &goodsDTO)
	if err != nil {
		//log.Errorf("get goods create error: %v", err.Error())
		return nil, errors.ToGrpcError(err)
	}
	ret.Id = goodsID
	return &ret, nil
}

func (gs *GoodsServer) DeleteGoods(ctx context.Context, info *goods_pb.DeleteGoodsInfo) (*empty.Empty, error) {
	err := gs.srv.Goods().Delete(ctx, uint64(info.Id))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) UpdateGoods(ctx context.Context, info *goods_pb.CreateGoodsInfo) (*empty.Empty, error) {
	goodsDO := do.GoodsDO{
		CategoryID:      info.CategoryId,
		BrandsID:        info.BrandId,
		OnSale:          info.OnSale,
		ShipFree:        info.ShipFree,
		IsNew:           info.IsNew,
		IsHot:           info.IsHot,
		Name:            info.Name,
		GoodsSn:         info.GoodsSn,
		MarketPrice:     info.MarketPrice,
		ShopPrice:       info.ShopPrice,
		GoodsBrief:      info.GoodsBrief,
		Images:          info.Images,
		DescImages:      info.DescImages,
		GoodsFrontImage: info.GoodsFrontImage,
	}
	goodsDTO := dto.GoodsDTO{GoodsDO: goodsDO}

	if err := gs.srv.Goods().Update(ctx, &goodsDTO); err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (gs *GoodsServer) GetGoodsDetail(ctx context.Context, req *goods_pb.GoodInfoRequest) (*goods_pb.GoodsInfoResponse, error) {
	good, err := gs.srv.Goods().Get(ctx, uint64(req.Id))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return DTOToResponse(good), nil
}
