package handler

import (
	"context"
	"go-learn/shop/shop_srvs/inventory_srv/global"
	"go-learn/shop/shop_srvs/inventory_srv/model"
	"go-learn/shop/shop_srvs/inventory_srv/proto"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 设置库存

	var inv model.Inventory

	// global.DB.First(&inv, req.GoodsId)
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)

	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)

	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

var m sync.Mutex

func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//   数据一致性 数据库事务
	tx := global.DB.Begin()
	m.Lock() // 加锁

	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减库存
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

	}

	tx.Commit() // 需要手动提交修改
	m.Unlock()  // 释放锁

	return &emptypb.Empty{}, nil
}

var m1 sync.Mutex

func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//   数据一致性 数据库事务
	tx := global.DB.Begin()
	m1.Lock() // 加锁

	// 库存归还 1.订单超时归还 2.订单创建失败 归还之前扣减  3.手动归还
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}
		// 退货
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}

	tx.Commit()
	m1.Unlock()

	return &emptypb.Empty{}, nil
}
