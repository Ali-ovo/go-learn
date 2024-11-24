package controller

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/app/shop_srv/goods/srv/internal/service"

	"github.com/golang/protobuf/ptypes/empty"
)

type GoodsServer struct {
	goods_pb.UnimplementedGoodsServer
	srv service.ServiceFactory
}

func (g GoodsServer) GetAllCategorysList(ctx context.Context, e *empty.Empty) (*goods_pb.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) GetCategorysList(ctx context.Context, request *goods_pb.CategoryListRequest) (*goods_pb.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) GetSubCategory(ctx context.Context, request *goods_pb.CategoryListRequest) (*goods_pb.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) CreateCategory(ctx context.Context, request *goods_pb.CategoryInfoRequest) (*goods_pb.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) DeleteCategory(ctx context.Context, request *goods_pb.DeleteCategoryRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) UpdateCategory(ctx context.Context, request *goods_pb.CategoryInfoRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) BrandList(ctx context.Context, request *goods_pb.BrandFilterRequest) (*goods_pb.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) CreateBrand(ctx context.Context, request *goods_pb.BrandRequest) (*goods_pb.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) DeleteBrand(ctx context.Context, request *goods_pb.BrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) UpdateBrand(ctx context.Context, request *goods_pb.BrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) BannerList(ctx context.Context, e *empty.Empty) (*goods_pb.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) CreateBanner(ctx context.Context, request *goods_pb.BannerRequest) (*goods_pb.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) DeleteBanner(ctx context.Context, request *goods_pb.BannerRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) UpdateBanner(ctx context.Context, request *goods_pb.BannerRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) CategoryBrandList(ctx context.Context, request *goods_pb.CategoryBrandFilterRequest) (*goods_pb.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) GetCategoryBrandList(ctx context.Context, request *goods_pb.CategoryInfoRequest) (*goods_pb.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) CreateCategoryBrand(ctx context.Context, request *goods_pb.CategoryBrandRequest) (*goods_pb.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) DeleteCategoryBrand(ctx context.Context, request *goods_pb.CategoryBrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) UpdateCategoryBrand(ctx context.Context, request *goods_pb.CategoryBrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodsServer) mustEmbedUnimplementedGoodsServer() {
	//TODO implement me
	panic("implement me")
}

func NewGoodsServer(srv service.ServiceFactory) goods_pb.GoodsServer {
	return &GoodsServer{srv: srv}
}

var _ goods_pb.GoodsServer = (*GoodsServer)(nil)

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
