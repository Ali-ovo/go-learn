package srv

import (
	"context"
	"shop/app/shop_srv/inventory/srv/internal/data/v1"
	"shop/app/shop_srv/inventory/srv/internal/domain/do"
	"shop/app/shop_srv/inventory/srv/internal/domain/dto"
	"shop/app/shop_srv/inventory/srv/internal/service"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/code"
	"sort"
	"strconv"
	"strings"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
)

const (
	inventoryLockPrefix = "inventory" // 锁商品库存用
	orderLockPrefix     = "order"     // 锁订单用
)

type inventoryService struct {
	data data.DataFactory
	pool redsyncredis.Pool // redis 池
}

func (is *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return is.data.Inventory().Create(ctx, nil, &inv.InventoryDO)
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

	var err error

	rs := redsync.New(is.pool)
	// 先按照商品的id排序, 然后从小到大逐个扣减库存, 防止锁竞争 和 防止死锁
	var detail = do.GoodsDetailList(details)
	sort.Sort(&detail)

	txn := is.data.Begin()
	defer func() { // 异常处理
		if p := recover(); p != nil {
			txn.Rollback()
			log.ErrorfC(ctx, "[回滚] 事务进行中出现异常: %v", p)
			return
		} else if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	sellDetail := do.StockSellDetailDO{
		OrderSn: ordersn,
		Status:  1,
		Detail:  detail,
	}

	mutexOrder := rs.NewMutex(strings.Join([]string{orderLockPrefix, ordersn}, "_"))
	if err = mutexOrder.Lock(); err != nil {
		log.InfofC(ctx, "订单 %s 获取锁失败", ordersn)
		return errors.WithCode(code.ErrRedisDatabase, err.Error())
	}
	defer mutexOrder.Unlock()

	for _, goodsInfo := range detail {
		mutexGoods := rs.NewMutex(strings.Join([]string{inventoryLockPrefix, strconv.Itoa(int(goodsInfo.Goods))}, "_"))
		if err = mutexGoods.Lock(); err != nil {
			log.InfofC(ctx, "商品 %d 获取锁失败", goodsInfo.Goods)
			return errors.WithCode(code.ErrRedisDatabase, err.Error())
		}
		defer mutexGoods.Unlock()

		// 查询库存信息是否存在
		var inv *do.InventoryDO
		inv, err = is.data.Inventory().Get(ctx, goodsInfo.Goods)
		if err != nil {
			log.Errorf("订单 %s 获取库存失败", ordersn)
			return err
		}

		// 判断库存是否充足
		if inv.Stocks < goodsInfo.Num {
			log.Errorf("商品 %d 库存 %d 不足, 现有库存 %d", goodsInfo.Goods, goodsInfo.Num, inv.Stocks)
			err = errors.WithCode(code.ErrInvNotEnough, "库存不足")
			return err
		}
		inv.Stocks -= goodsInfo.Num

		result := is.data.Inventory().Reduce(ctx, txn, goodsInfo.Goods, goodsInfo.Num)
		if result.RowsAffected == 0 || result.Error != nil {
			log.Errorf("订单 %s 扣减库存失败", ordersn)
			return err
		}
	}

	err = is.data.Inventory().CreateStockSellDetail(ctx, txn, &sellDetail)
	if err != nil {
		log.Errorf("订单 %s 创建扣减库存记录失败", ordersn)
		return err
	}

	txn.Commit()
	return nil
}

func (is *inventoryService) Repack(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单 %s 归还库存", ordersn)

	var err error

	rs := redsync.New(is.pool)

	txn := is.data.Begin()
	defer func() { // 异常处理
		if p := recover(); p != nil {
			txn.Rollback()
			log.ErrorfC(ctx, "[回滚]事务进行中出现异常: %v", p)
			return
		} else if err != nil {
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	// 主动取消 网络问题引起的重试 超时取消 退款取消
	mutexOrder := rs.NewMutex(strings.Join([]string{orderLockPrefix, ordersn}, "_"))
	if err = mutexOrder.Lock(); err != nil {
		log.InfofC(ctx, "订单 %s 获取锁失败", ordersn)
		return errors.WithCode(code.ErrRedisDatabase, err.Error())
	}
	defer mutexOrder.Unlock()

	// 获取订单
	sellDetail, err := is.data.Inventory().GetSellDetail(ctx, txn, ordersn)
	if err != nil {
		if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
			log.Errorf("[忽略]订单 %s 扣减库存记录不存在", ordersn)
			return nil
		}
		log.Errorf("订单 %s 获取扣减库存记录失败", ordersn)
		return err
	}

	if sellDetail.Status == 2 {
		log.Infof("订单 %s 扣减库存记录已经归还, 忽略", ordersn)
		return nil
	}

	var detail = do.GoodsDetailList(details)
	sort.Sort(&detail)

	for _, goodsInfo := range sellDetail.Detail {
		mutexGoods := rs.NewMutex(strings.Join([]string{inventoryLockPrefix, strconv.Itoa(int(goodsInfo.Goods))}, "_"))
		if err = mutexGoods.Lock(); err != nil {
			log.InfofC(ctx, "订单 %s 获取锁失败", ordersn)
			return errors.WithCode(code.ErrRedisDatabase, err.Error())
		}
		defer mutexGoods.Unlock()

		inv, err := is.data.Inventory().Get(ctx, goodsInfo.Goods)
		if err != nil {
			log.Errorf("订单 %s 获取库存失败", ordersn)
			return err
		}
		inv.Stocks += goodsInfo.Num

		err = is.data.Inventory().Increase(ctx, txn, goodsInfo.Goods, goodsInfo.Num)
		if err != nil {
			log.Errorf("订单 %s 归还库存失败", ordersn)
			return err
		}
	}
	err = is.data.Inventory().UpdateStockSellDetailStatus(ctx, txn, ordersn, 2)
	if err != nil {
		log.Errorf("订单 %s 更新扣减库存记录失败", ordersn)
		return err
	}
	txn.Commit()
	return nil
}

func newInventory(srv *serviceFactory) service.InventorySrv {
	return &inventoryService{
		data: srv.data,
		pool: srv.pool,
	}
}
