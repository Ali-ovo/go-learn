package srv

import (
	"context"
	"encoding/json"
	"fmt"
	inventory_pb "shop/api/inventory/v1"
	order_pb "shop/api/order/v1"
	"shop/app/shop_srv/order/srv/config"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/pkg/log"
	"shop/pkg/options"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/dtm-labs/client/dtmgrpc"
)

func RocketmqConsumer(opts *config.Config) {
	var err error

	rocketMQ := conn.RocketMQ{
		Addr:          opts.Rocketmq.Addr,
		ConsumerGroup: opts.Rocketmq.GroupName,
		Retry:         opts.Rocketmq.Retry,
		Consumer: []conn.Consumer{
			{
				Topic:    "order_timeout",
				Selector: consumer.MessageSelector{},
				Func:     OrderTimeout(opts.Dtm), // TODO 回滚逻辑
			},
		},
	}

	// 连接 rockerMq 启动消费者消费消息
	err = rocketMQ.ConsumerStart()
	if err != nil {
		panic(err)
	}
	// 连接 rockerMq 启动生产者客户端
	err = rocketMQ.ProducerStart()
	if err != nil {
		panic(err)
	}

	log.Info("[order-srv] 启动 Rocketmq 完成")
}

func OrderTimeout(opt *options.DtmOptions) func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		var err error

		for i := range msgs {
			var orderInfo dto.OrderDTO
			_ = json.Unmarshal(msgs[i].Body, &orderInfo)

			var goodsInfo []*inventory_pb.GoodsInvInfo
			var orderItems []*order_pb.OrderItemResponse
			for _, item := range orderInfo.OrderGoods {
				goodsInfo = append(goodsInfo, &inventory_pb.GoodsInvInfo{
					GoodsId: item.Goods,
					Num:     item.Nums,
				})

				orderItems = append(orderItems, &order_pb.OrderItemResponse{
					GoodsId: item.Goods,
					Num:     item.Nums,
				})
			}

			iReq := &inventory_pb.SellInfo{
				GoodsInfo: goodsInfo,
				OrderSn:   orderInfo.OrderSn,
			}
			oReq := &order_pb.OrderRequest{
				OrderSn:    orderInfo.OrderSn,
				UserId:     orderInfo.User,
				Address:    orderInfo.Address,
				Name:       orderInfo.SignerName,
				Mobile:     orderInfo.SingerMobile,
				Post:       orderInfo.Post,
				OrderItems: orderItems, // 订单商品
			}

			// 使用 dtm 进行分布式归还
			saga := dtmgrpc.NewSagaGrpc(opt.GrpcServer, "rollback_"+orderInfo.OrderSn).
				Add(opt.AccessPath["inventory"]+"/Inventory/Reback", opt.AccessPath["inventory"]+"/Inventory/Sell", iReq).   // 回退库存
				Add(opt.AccessPath["order"]+"/Order/RollBackOrder", opt.AccessPath["order"]+"/Order/RollBackOrderCom", oReq) // 删除订单  // TODO 恢复购物车商品

			saga.WaitResult = true // 设置: 等待执行完成
			// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
			if err = saga.Submit(); err != nil {
				fmt.Println(err)
				return consumer.ConsumeRetryLater, nil // 如果有错误 重新消费
			}
		}

		fmt.Println(msgs)
		return consumer.ConsumeSuccess, nil
	}
}
