package srv

import (
	"context"
	"database/sql"
	"encoding/json"
	goods_pb "shop/api/goods/v1"
	inventory_pb "shop/api/inventory/v1"
	order_pb "shop/api/order/v1"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	"shop/app/shop_srv/order/srv/internal/service/v1"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	code2 "shop/pkg/code"
	"shop/pkg/options"

	"github.com/dtm-labs/client/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type orderService struct {
	data    data.DataFactory
	dtmOpts *options.DtmOptions
}

func (os *orderService) Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error) {
	order, err := os.data.Orders().Get(ctx, orderSn)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrOrderNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &dto.OrderDTO{OrderInfoDO: *order}, nil
}

func (os *orderService) List(ctx context.Context, userID int64, meta metav1.ListMeta, orderby []string) (*dto.OrderDTOList, error) {
	var ret dto.OrderDTOList

	ordersList, err := os.data.Orders().List(ctx, userID, meta, orderby)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.TotalCount = ordersList.TotalCount
	for _, value := range ordersList.Items {
		ret.Items = append(ret.Items, &dto.OrderDTO{
			OrderInfoDO: *value,
		})
	}
	return &ret, nil
}

func (os *orderService) Submit(ctx context.Context, order *dto.OrderDTO) error {
	var err error

	// 先从购物车中获取商品信息
	list, err := os.data.ShopCarts().List(ctx, order.User, true, metav1.ListMeta{}, []string{})
	if err != nil {
		log.Errorf("获取购物车信息失败, err:%v", err)
		return err
	}
	if len(list.Items) == 0 {
		log.Errorf("获取购物车没有商品, 无法下单")
		return errors.WithCode(code2.ErrNotGoodsSelect, "购物车中没有选择商品")
	}

	var orderGoods []*do.OrderGoods
	var orderItems []*order_pb.OrderItemResponse
	for _, value := range list.Items {
		orderGoods = append(orderGoods, &do.OrderGoods{
			Goods: value.Goods,
			Nums:  value.Nums,
		})

		orderItems = append(orderItems, &order_pb.OrderItemResponse{
			GoodsId: value.Goods,
			Num:     value.Nums,
		})
	}
	// 查询 真实购物车商品(因为前端传递可能出现问题, 我们可以自己去数据库中获取)
	order.OrderGoods = orderGoods

	// 基于可靠消息最终一致性的思想, saga事务来解决订单生成的问题
	var goodsInfo []*inventory_pb.GoodsInvInfo
	for _, item := range order.OrderGoods {
		goodsInfo = append(goodsInfo, &inventory_pb.GoodsInvInfo{
			GoodsId: item.Goods,
			Num:     item.Nums,
		})
	}
	iReq := &inventory_pb.SellInfo{
		GoodsInfo: goodsInfo,
		OrderSn:   order.OrderSn,
	}
	oReq := &order_pb.OrderRequest{
		OrderSn:    order.OrderSn,
		UserId:     order.User,
		Address:    order.Address,
		Name:       order.SignerName,
		Mobile:     order.SingerMobile,
		Post:       order.Post,
		OrderItems: orderItems, // 订单商品
	}

	saga := dtmgrpc.NewSagaGrpc(os.dtmOpts.GrpcServer, "create"+order.OrderSn).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 逆向操作为url: qsBusi+"/TransOutCom"
		Add(os.dtmOpts.AccessPath["inventory"]+"/Inventory/Sell", os.dtmOpts.AccessPath["inventory"]+"/Inventory/Reback", iReq). // 扣减库存
		Add(os.dtmOpts.AccessPath["order"]+"/Order/CreateOrder", os.dtmOpts.AccessPath["order"]+"/Order/CreateOrderCom", oReq)   // 创建订单
	saga.WaitResult = true // 设置: 等待执行完成
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	if err = saga.Submit(); err != nil {
		return errors.WithCode(code2.ErrOrderDtm, "Dtm 事务提交失败")
	}
	// TODO 查询 dtm gid 查看状态 返回对应错误信息
	return err
}

func (os *orderService) Create(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1. 生成 orderinfo  订单 表数据
		2. 生成 ordergoods 订单商品 表数据
		3. 根据 order 找到对应的购物车条目, 删除购物车条目
	*/

	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	txn := os.data.Begin()
	sourceTx := txn.Statement.ConnPool.(*sql.Tx)

	err := barrier.Call(sourceTx, func(tx *sql.Tx) error {
		var goodsIds []int64
		for _, value := range order.OrderGoods {
			goodsIds = append(goodsIds, value.Goods)
		}

		//// 测试
		//return status.Error(codes.Aborted, "create order failed")

		goods, err := os.data.Goods().BatchGetGoods(ctx, &goods_pb.BatchGoodsIdInfo{Id: goodsIds})
		if err != nil {
			log.Errorf("批量获取商品信息失败, goodsIds: %v, err: %v", goodsIds, err)
			return err // 查询失败 重复查询
		}
		if len(goods.Data) != len(goodsIds) {
			log.Errorf("批量获取商品信息失败, goodsIds: %v, 返回值: %v, err: %v", goodsIds, goods.Data, err)
			return status.Error(codes.Aborted, "商品不存在 or 部分商品不存在") // 回滚
		}

		var goodsMap = make(map[int64]*goods_pb.GoodsInfoResponse)
		for _, value := range goods.Data {
			goodsMap[value.Id] = value
		}
		var orderAmount float32 // 订单总价
		for _, value := range order.OrderGoods {
			orderAmount += goodsMap[value.Goods].ShopPrice * float32(value.Nums)
			value.GoodsName = goodsMap[value.Goods].Name
			value.GoodsPrice = goodsMap[value.Goods].ShopPrice
			value.GoodsImage = goodsMap[value.Goods].GoodsFrontImage
		}

		if result := os.data.Orders().Create(ctx, txn, &order.OrderInfoDO); result.RowsAffected == 0 {
			log.Errorf("创建订单失败, err:%v", err)
			if result.Error != nil {
				return status.Error(codes.FailedPrecondition, result.Error.Error()) // 重试
			}
			return status.Error(codes.Aborted, "创建订单失败") // 回滚
		}

		if result := os.data.ShopCarts().DeleteByGoodsIDs(ctx, txn, order.User, goodsIds); result.RowsAffected == 0 {
			log.Errorf("删除购物车商品失败, goodsIds:%v, err:%v", goodsIds, result.Error)
			if result.Error != nil {
				return status.Error(codes.FailedPrecondition, result.Error.Error())
			}
			return status.Error(codes.Aborted, "删除购物车商品失败") // 回滚
		}

		jsonByte, _ := json.Marshal(order)

		// TODO 创建二进制 数据
		err = conn.Message("order_timeout", jsonByte, 15)
		if err != nil {
			return status.Error(codes.Aborted, "发送延迟消息 订单超时失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (os *orderService) CreateCom(ctx context.Context, order *dto.OrderDTO) error {
	/*
		1. 删除 orderinfo 表数据
		2. 删除 ordergoods 表
		3. 根据 order 找到对应的购物车条目, 添加购物车条目

		PS： 不需要补偿内容
		因为 如果是 mysql的错误, 会重试
		如果 是 其他错误, 会进行 mysql 事务回滚, 数据未被更改,
		所以不需要额外的补偿

		如果真需要回滚
		应先查询订单身份存在, 如果存在 删除相关记录 同时 回滚 购物车商品记录
	*/
	return nil
}

func (os *orderService) RollBack(ctx context.Context, order *dto.OrderDTO) error {
	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	txn := os.data.Begin()
	sourceTx := txn.Statement.ConnPool.(*sql.Tx)

	err := barrier.Call(sourceTx, func(tx *sql.Tx) error {
		_, err := os.data.Orders().Get(ctx, order.OrderSn)
		if err != nil {
			log.Errorf("查询订单失败, err:%v", err)
			if err != nil {
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}
			return status.Error(codes.Aborted, "无此订单") // PS: 理论上不可能需要回滚
		}

		if result := os.data.Orders().Update(ctx, txn, &do.OrderInfoDO{
			OrderSn: order.OrderSn,
			Status:  "TRADE_CLOSED",
		}); result.RowsAffected == 0 {
			if result.Error != nil {
				return status.Error(codes.FailedPrecondition, err.Error()) // 重试
			}
			return status.Error(codes.Aborted, "订单更新失败") // PS: 理论上不可能需要回滚
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (os *orderService) RollBackCom(ctx context.Context, order *dto.OrderDTO) error {
	return nil
}

func (os *orderService) Update(ctx context.Context, order *dto.OrderDTO) error {
	//TODO implement me
	panic("implement me")
}

func newOrders(srv *serviceFactory) service.OrderSrv {
	return &orderService{
		data:    srv.data,
		dtmOpts: srv.dtmOpts,
	}
}

// TODO 订单延迟后 使用 dtm 进行库存归还 购物车商品归还 订单状态修改
