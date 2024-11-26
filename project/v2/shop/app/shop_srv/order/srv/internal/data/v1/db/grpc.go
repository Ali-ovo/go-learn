package db

import (
	"context"
	goods_pb "shop/api/goods/v1"
	inventory_pb "shop/api/inventory/v1"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver"
	"shop/pkg/code"
)

const goodsServiceName = "discovery:///shop-goods-srv"
const InventoryServiceName = "discovery:///shop-Inventory-srv"

func NewGoodsServiceClient(r registry.Discovery) (goods_pb.GoodsClient, error) {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		// 设置负载均衡
		rpcserver.WithBanlancerName("selector"),
		rpcserver.WithDiscovery(r),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpcserver.WithEndpoint(goodsServiceName),
		//rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),		// 这是自己封装的 链路追踪
		rpcserver.WithClientEnableTracing(true),
		//rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectGRPC, "failed to get grpc store factory")
	}

	c := goods_pb.NewGoodsClient(conn)
	return c, nil
}

func NewInventoryServiceClient(r registry.Discovery) (inventory_pb.InventoryClient, error) {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		// 设置负载均衡
		rpcserver.WithBanlancerName("selector"),
		rpcserver.WithDiscovery(r),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpcserver.WithEndpoint(InventoryServiceName),
		//rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),		// 这是自己封装的 链路追踪
		rpcserver.WithClientEnableTracing(true),
		//rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		return nil, errors.WithCode(code.ErrConnectGRPC, "failed to get grpc store factory")
	}

	c := inventory_pb.NewInventoryClient(conn)
	return c, nil
}
