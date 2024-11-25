package srv

import (
	"context"
	"database/sql"
	"shop/app/shop_srv/inventory/srv/internal/data/v2"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/app/shop_srv/inventory/srv/internal/domain/dto"
	"shop/app/shop_srv/inventory/srv/internal/service"
	code2 "shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/code"
	"sort"
	"strconv"
	"strings"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	inventoryLockPrefix = "inventory" // 锁商品库存用
)

type inventoryService struct {
	data data.DataFactory
	pool redsyncredis.Pool // redis 池
}

func (is *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return is.data.Inventory().Create(ctx, &inv.InventoryDO)
}

func (is *inventoryService) Get(ctx context.Context, goodsID int64) (*dto.InventoryDTO, error) {
	inv, err := is.data.Inventory().Get(ctx, goodsID)
	if err != nil {
		return nil, err
	}
	return &dto.InventoryDTO{InventoryDO: *inv}, nil
}

func (is *inventoryService) Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单 %s 扣减库存", ordersn)

	rs := redsync.New(is.pool)

	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	txn := is.data.Begin()
	sourceTx := txn.Statement.ConnPool.(*sql.Tx)

	err := barrier.Call(sourceTx, func(tx *sql.Tx) error {
		// 先按照商品的id排序, 然后从小到大逐个扣减库存, 防止锁竞争 和 防止死锁
		var detail = do.GoodsDetailList(details)
		sort.Sort(&detail)

		sellDetail := do.StockSellDetailDO{
			OrderSn: ordersn,
			Status:  1,
			Detail:  detail,
		}

		for _, goodsInfo := range detail {
			mutexGoods := rs.NewMutex(strings.Join([]string{inventoryLockPrefix, strconv.Itoa(int(goodsInfo.Goods))}, "_"))
			if err := mutexGoods.Lock(); err != nil {
				log.InfofC(ctx, "商品 %d 获取锁失败", goodsInfo.Goods)
				return status.Error(codes.Aborted, err.Error()) // 回滚

			}
			defer mutexGoods.Unlock()

			// 查询库存信息是否存在
			var inv *do.InventoryDO
			inv, err := is.data.Inventory().Get(ctx, goodsInfo.Goods)
			if err != nil {
				log.Errorf("订单 %s 获取库存失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}

			// 判断库存是否充足
			if inv.Stocks < goodsInfo.Num {
				log.Errorf("商品 %d 库存 %d 不足, 现有库存 %d", goodsInfo.Goods, goodsInfo.Num, inv.Stocks)
				return status.Error(codes.Aborted, "库存不足") // 回滚
			}
			inv.Stocks -= goodsInfo.Num

			result, err := is.data.Inventory().Reduce(ctx, sourceTx, goodsInfo.Goods, goodsInfo.Num)
			if err != nil {
				log.Errorf("订单 %s 扣减库存失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}
			if rows, _ := result.RowsAffected(); rows == 0 {
				return status.Error(codes.Aborted, "查询不到商品库存信息") // 回滚
			}
		}

		_, err := is.data.Inventory().CreateStockSellDetail(ctx, sourceTx, &sellDetail)
		if err != nil {
			if errors.Code(err) != code2.ErrDecodingJSON {
				log.Errorf("订单 %s 创建扣减库存记录失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 数据库:重试
			} else {
				log.Errorf("订单 %s JSON 映射失败回滚", ordersn)
				return status.Error(codes.Aborted, err.Error()) // json:回滚
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (is *inventoryService) Repack(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单 %s 归还库存", ordersn)

	rs := redsync.New(is.pool)

	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	txn := is.data.Begin()
	sourceTx := txn.Statement.ConnPool.(*sql.Tx)

	err := barrier.Call(sourceTx, func(tx *sql.Tx) error {
		// 获取订单
		sellDetail, err := is.data.Inventory().GetSellDetail(ctx, ordersn)
		if err != nil {
			if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
				log.Errorf("[忽略]订单 %s 扣减库存记录不存在", ordersn)
				return nil // 订单不存在 说明还没下单 忽略 (理论上不存在这个问题)
			}
			log.Errorf("订单 %s 获取扣减库存记录失败", ordersn)
			return status.Error(codes.FailedPrecondition, err.Error()) // 重试 可能 mysql 出现问题
		}

		if sellDetail.Status == 2 {
			log.Infof("订单 %s 扣减库存记录已经归还, 忽略", ordersn)
			return nil // 已经归还 忽略
		}

		var detail = do.GoodsDetailList(details)
		sort.Sort(&detail)

		for _, goodsInfo := range sellDetail.Detail {
			mutexGoods := rs.NewMutex(strings.Join([]string{inventoryLockPrefix, strconv.Itoa(int(goodsInfo.Goods))}, "_"))
			if err = mutexGoods.Lock(); err != nil {
				log.InfofC(ctx, "订单 %s 获取锁失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试 redis 出现问题
			}
			defer mutexGoods.Unlock()

			inv, err := is.data.Inventory().Get(ctx, goodsInfo.Goods)
			if err != nil {
				log.Errorf("订单 %s 获取商品库存失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}
			inv.Stocks += goodsInfo.Num

			result, err := is.data.Inventory().Increase(ctx, sourceTx, goodsInfo.Goods, goodsInfo.Num)
			if err != nil {
				log.Errorf("订单 %s 归还库存失败", ordersn)
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}
			if rows, _ := result.RowsAffected(); rows == 0 {
				return status.Error(codes.Aborted, "查询不到商品库存信息") // 回滚
			}
		}
		result, err := is.data.Inventory().UpdateStockSellDetailStatus(ctx, sourceTx, ordersn, 2)
		if err != nil {
			log.Errorf("订单 %s 更新订单状态失败", ordersn)
			return status.Error(codes.FailedPrecondition, err.Error()) // 重试
		}
		if rows, _ := result.RowsAffected(); rows == 0 {
			return status.Error(codes.Aborted, "查询不到此订单信息") // 回滚
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func newInventory(srv *serviceFactory) service.InventorySrv {
	return &inventoryService{
		data: srv.data,
		pool: srv.pool,
	}
}
