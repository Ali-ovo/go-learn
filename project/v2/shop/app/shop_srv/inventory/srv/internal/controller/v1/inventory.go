package controller

import (
	"context"
	inventory_pb "shop/api/inventory/v1"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/app/shop_srv/inventory/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SetInv
//
//	@Description: 设置库存
//	@receiver is
//	@param ctx
//	@param info
//	@return *empty.Empty
//	@return error
func (is *InventoryServer) SetInv(ctx context.Context, info *inventory_pb.GoodsInvInfo) (*empty.Empty, error) {
	invDTO := &dto.InventoryDTO{}
	invDTO.Goods = info.GoodsId
	invDTO.Stocks = info.Num
	err := is.srv.Inventory().Create(ctx, invDTO)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// InvDetail
//
//	@Description: 根据商品的id查询库存
//	@receiver is
//	@param ctx
//	@param info
//	@return *inventory_pb.GoodsInvInfo
//	@return error
func (is *InventoryServer) InvDetail(ctx context.Context, info *inventory_pb.GoodsInvInfo) (*inventory_pb.GoodsInvInfo, error) {
	inv, err := is.srv.Inventory().Get(ctx, info.GoodsId)
	if err != nil {
		return nil, err
	}
	return &inventory_pb.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// Sell
//
//	@Description: 扣减库存
//	@receiver is
//	@param ctx
//	@param info
//	@return *empty.Empty
//	@return error
func (is *InventoryServer) Sell(ctx context.Context, info *inventory_pb.SellInfo) (*empty.Empty, error) {
	var detail []do.GoodsDetail
	for _, value := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{Goods: value.GoodsId, Num: value.Num})
	}
	err := is.srv.Inventory().Sell(ctx, info.OrderSn, detail)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	//time.Sleep(3 * time.Second)
	//return nil, status.Error(codes.Aborted, "测试")
	return &emptypb.Empty{}, nil
}

// Reback
//
//	@Description: 归还库存
//	@receiver is
//	@param ctx
//	@param info
//	@return *empty.Empty
//	@return error
func (is *InventoryServer) Reback(ctx context.Context, info *inventory_pb.SellInfo) (*empty.Empty, error) {
	log.Infof("订单 %s 归还库存", info.OrderSn)
	var detail []do.GoodsDetail
	for _, v := range info.GoodsInfo {
		detail = append(detail, do.GoodsDetail{Goods: v.GoodsId, Num: v.Num})
	}
	err := is.srv.Inventory().Repack(ctx, info.OrderSn, detail)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &emptypb.Empty{}, nil
}
