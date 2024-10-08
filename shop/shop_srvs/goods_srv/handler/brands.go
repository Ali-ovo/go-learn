package handler

import (
	"context"
	"go-learn/shop/shop_srvs/goods_srv/global"
	"go-learn/shop/shop_srvs/goods_srv/model"
	"go-learn/shop/shop_srvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {

	var brands []model.Brands

	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)

	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	var brandResponse []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponse = append(brandResponse, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return &proto.BrandListResponse{
		Data:  brandResponse,
		Total: int32(total),
	}, nil
}

// 新建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {

	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Create(brand)

	return &proto.BrandInfoResponse{
			Id: int32(brand.ID),
		},
		nil

}

// 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {

	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}

	if result := global.DB.First(&brands, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}

	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	// 保存更新后的品牌信息
	if result := global.DB.Save(&brands); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败")
	}

	return &emptypb.Empty{}, nil
}
