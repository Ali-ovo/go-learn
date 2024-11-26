package controller

import (
	"context"
	order_pb "shop/api/order/v1"

	"github.com/golang/protobuf/ptypes/empty"
)

func (os *OrderServer) CartItemList(ctx context.Context, info *order_pb.UserInfo) (*order_pb.CartItemListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *OrderServer) CreateCartItem(ctx context.Context, request *order_pb.CartItemRequest) (*order_pb.ShopCartInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (os *OrderServer) UpdateCartItem(ctx context.Context, request *order_pb.CartItemRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (os *OrderServer) DeleteCartItem(ctx context.Context, request *order_pb.CartItemRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}
