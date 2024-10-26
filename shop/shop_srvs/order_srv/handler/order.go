package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-learn/shop/shop_srvs/order_srv/global"
	"go-learn/shop/shop_srvs/order_srv/model"
	"go-learn/shop/shop_srvs/order_srv/proto"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// 订单
var p rocketmq.TransactionProducer

var globalOrderListener *OrderListener

var p1 rocketmq.Producer

func initP1() error {
	var err error
	p1, err = rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.189.128:9876"}),
		producer.WithGroupName("timeout_order"),
	)

	if err != nil {
		zap.S().Errorf("生成 producer 失败: %s", err.Error())
		return err
	}

	return p1.Start()
}

func initProducer(orderListener *OrderListener) error {
	var err error
	p, err = rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNameServer([]string{"192.168.189.128:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成 producer 失败: %s", err.Error())
		return err
	}
	return p.Start()
}

func GenerateOrderSn(userId int32) string {
	// 生成订单号
	// 年月日时分秒+用户ID+2位随机数

	now := time.Now()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), r.Intn(90)+10)

	return orderSn
}

// 购物车
func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart

	var rsp proto.CartItemListResponse
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		// 如果记录存在，则更新数量
		shopCart.Nums += req.Nums
	} else {
		// 如果记录不存在，则创建
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)

	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}

	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	return &emptypb.Empty{}, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "购物车为空"
		return primitive.RollbackMessageState
	}

	for _, shopCarts := range shopCarts {
		goodsIds = append(goodsIds, shopCarts.Goods)
		goodsNumsMap[shopCarts.Goods] = shopCarts.Nums
	}

	// 商品微服务 gin 调用
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "获取商品信息失败"
		return primitive.RollbackMessageState
	}

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])

		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})

	}

	// 库存微服务
	if _, err := global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfo,
		OrderSn:   orderInfo.OrderSn,
	}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	// 生成订单表
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount

	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()

		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}
	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGoods := range orderGoods {
		orderGoods.Order = orderInfo.ID // 订单号
	}

	// 生成订单商品表
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()

		o.Code = codes.Internal
		o.Detail = "创建订单商品失败"
		return primitive.CommitMessageState
	}

	// 清空购物车
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()

		o.Code = codes.Internal
		o.Detail = "清空购物车失败"
		return primitive.CommitMessageState
	}

	// 发送延时消息
	if p1 == nil {
		if err := initP1(); err != nil {
			zap.S().Errorf("初始化生产者失败: %s", err.Error())
			return primitive.CommitMessageState
		}
	}

	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = p1.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	if globalOrderListener == nil {
		globalOrderListener = &OrderListener{Ctx: ctx}
	}

	// orderListener := OrderListener{Ctx: ctx}
	// p, err := rocketmq.NewTransactionProducer(
	// 	&orderListener,
	// 	producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.189.128:9876"})),
	// )

	// if err != nil {
	// 	zap.S().Errorf("生成 producer 失败: %s", err.Error())
	// 	return nil, err

	// }

	// if err = p.Start(); err != nil {
	// 	zap.S().Errorf("启动 producer 失败: %s", err.Error())
	// 	return nil, err
	// }

	if p == nil {
		if err := initProducer(globalOrderListener); err != nil {
			zap.S().Errorf("初始化生产者失败: %s", err.Error())
			return nil, err
		}
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	jsonString, _ := json.Marshal(order)

	_, err := p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage("order_reback", jsonString))

	if err != nil {
		fmt.Printf("发送消息失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}

	if globalOrderListener.Code != codes.OK {
		return nil, status.Error(globalOrderListener.Code, globalOrderListener.Detail)
	}

	return &proto.OrderInfoResponse{
		Id:      globalOrderListener.ID,
		OrderSn: order.OrderSn,
		Total:   globalOrderListener.OrderAmount,
	}, nil
}

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)

	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		//查询订单的支付状态，如果已支付什么都不做，如果未支付，归还库存
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			if p1 == nil {
				if err := initP1(); err != nil {
					zap.S().Errorf("初始化生产者失败: %s", err.Error())
					return consumer.ConsumeRetryLater, nil
				}
			}

			_, err := p1.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}

			return consumer.ConsumeSuccess, nil
		}
	}
	return consumer.ConsumeSuccess, nil
}
