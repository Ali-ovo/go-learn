package handler

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/inventory_srv/global"
	"go-learn/shop/shop_srvs/inventory_srv/model"
	"go-learn/shop/shop_srvs/inventory_srv/proto"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"

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
	var inv model.InventoryNew
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks - inv.Freeze,
	}, nil
}

// var m sync.Mutex

func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	//   数据一致性 数据库事务
	tx := global.DB.Begin()
	// m.Lock() // 加锁

	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		// 	tx.Rollback() // 回滚之前的操作
		// 	return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		// }

		// for {

		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))

		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis锁失败")
		}

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

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis锁失败")
		}

		// 乐观锁
		// if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
		// 	zap.S().Info("库存扣减失败")
		// } else {
		// 	break
		// }
		// }
	}

	tx.Commit() // 需要手动提交修改
	// m.Unlock()  // 释放锁

	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) TrySell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	//   数据一致性 数据库事务
	tx := global.DB.Begin()

	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew

		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))

		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis锁失败")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减库存
		inv.Freeze += goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis锁失败")
		}

	}

	tx.Commit() // 需要手动提交修改

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

func (*InventoryServer) ConfirmSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew

		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= goodInfo.Num
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}

	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) CancelSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew

		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}

	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}
