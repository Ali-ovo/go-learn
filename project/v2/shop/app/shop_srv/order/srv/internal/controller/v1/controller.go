package controller

import (
	order_pb "shop/api/order/v1"
	"shop/app/shop_srv/order/srv/internal/service/v1"
)

type OrderServer struct {
	order_pb.UnimplementedOrderServer
	srv service.ServiceFactory
}

func NewOrderServer(srv service.ServiceFactory) order_pb.OrderServer {
	return &OrderServer{srv: srv}
}
