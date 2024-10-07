package handler

import (
	"go-learn/shop/shop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// 商品接口
// func (s *GoodsServer) GoodsList(ctx context.Context, in *proto.GoodsFilterRequest, opts ...grpc.CallOption) (*proto.GoodsListResponse, error) {

// 	return nil, nil
// }

// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
// BatchGetGoods(ctx context.Context, in *BatchGoodsIdInfo, opts ...grpc.CallOption) (*GoodsListResponse, error)
// CreateGoods(ctx context.Context, in *CreateGoodsInfo, opts ...grpc.CallOption) (*GoodsInfoResponse, error)
// DeleteGoods(ctx context.Context, in *DeleteGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
// UpdateGoods(ctx context.Context, in *CreateGoodsInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
// GetGoodsDetail(ctx context.Context, in *GoodInfoRequest, opts ...grpc.CallOption) (*GoodsInfoResponse, error)
