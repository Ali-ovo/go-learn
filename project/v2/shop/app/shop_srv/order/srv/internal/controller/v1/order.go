package controller

import (
	"context"
	order_pb "shop/api/order/v1"
	"shop/app/shop_srv/order/srv/internal/domain/do"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
)

/*
	订单提交的时候应该是先 生成订单号
	生成订单号要单独做一个接口，订单查询，以及一系列的关联我们应该采用order_sn，不要再去采用id去关联
	PS 为了幂等性 如果在srv层 生成 订单号, 重复发送创建订单 会导致重复创建
*/

func (os *OrderServer) SubmitOrder(ctx context.Context, request *order_pb.OrderRequest) (*empty.Empty, error) {
	//从购物车中得到选中的商品
	orderDTO := dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			User:         request.UserId,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
			OrderSn:      request.OrderSn,
		},
	}
	err := os.srv.Orders().Submit(ctx, &orderDTO)
	if err != nil {
		log.Errorf("新建订单失败: %v", err)
		return nil, errors.ToGrpcError(err)
	}

	return &empty.Empty{}, nil
}

// CreateOrder 这个是给分布式事务 saga 调用, 目前没有为 api 提供
func (os *OrderServer) CreateOrder(ctx context.Context, request *order_pb.OrderRequest) (*empty.Empty, error) {
	orderGoods := make([]*do.OrderGoods, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoods{
			Goods: item.GoodsId,
			Nums:  item.Num,
		}
	}

	err := os.srv.Orders().Create(ctx, &dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			OrderGoods:   orderGoods,
			User:         request.UserId,
			OrderSn:      request.OrderSn,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
		},
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (os *OrderServer) CreateOrderCom(ctx context.Context, request *order_pb.OrderRequest) (*empty.Empty, error) {
	orderGoods := make([]*do.OrderGoods, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoods{
			Goods: item.GoodsId,
			Nums:  item.Num,
		}
	}

	err := os.srv.Orders().CreateCom(ctx, &dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			OrderGoods:   orderGoods,
			User:         request.UserId,
			OrderSn:      request.OrderSn,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
		},
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (os *OrderServer) RollBackOrder(ctx context.Context, request *order_pb.OrderRequest) (*empty.Empty, error) {
	orderGoods := make([]*do.OrderGoods, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoods{
			Goods: item.GoodsId,
			Nums:  item.Num,
		}
	}

	err := os.srv.Orders().RollBack(ctx, &dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			OrderGoods:   orderGoods,
			User:         request.UserId,
			OrderSn:      request.OrderSn,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
		},
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (os *OrderServer) RollBackOrderCom(ctx context.Context, request *order_pb.OrderRequest) (*empty.Empty, error) {
	orderGoods := make([]*do.OrderGoods, len(request.OrderItems))
	for i, item := range request.OrderItems {
		orderGoods[i] = &do.OrderGoods{
			Goods: item.GoodsId,
			Nums:  item.Num,
		}
	}

	err := os.srv.Orders().RollBackCom(ctx, &dto.OrderDTO{
		OrderInfoDO: do.OrderInfoDO{
			OrderGoods:   orderGoods,
			User:         request.UserId,
			OrderSn:      request.OrderSn,
			Address:      request.Address,
			SignerName:   request.Name,
			SingerMobile: request.Mobile,
			Post:         request.Post,
		},
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}

func (os *OrderServer) OrderList(ctx context.Context, request *order_pb.OrderFilterRequest) (*order_pb.OrderListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *OrderServer) OrderDetail(ctx context.Context, request *order_pb.OrderRequest) (*order_pb.OrderInfoDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *OrderServer) UpdateOrderStatus(ctx context.Context, status *order_pb.OrderStatus) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}
