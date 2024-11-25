package controller

import (
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/internal_srv/domain/dto"
	"shop/app/shop/api/internal_srv/service"
)

type GoodsServer struct {
	goods_pb.UnimplementedGoodsServer
	srv service.ServiceFactory
}

func NewGoodsServer(srv service.ServiceFactory) goods_pb.GoodsServer {
	return &GoodsServer{srv: srv}
}

func DTOToResponse(goods *dto.GoodsDTO) *goods_pb.GoodsInfoResponse {
	goodsInfoRsp := goods_pb.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &goods_pb.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &goods_pb.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}

	return &goodsInfoRsp
}
