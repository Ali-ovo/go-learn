package goods

import (
	"context"
	goods_pb "shop/api/goods/v1"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver"
)

const goodsserviceName = "discovery:///shop-goods-srv"

func NewGoodsServiceClient(r registry.Discovery) goods_pb.GoodsClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		// 设置负载均衡
		rpcserver.WithBanlancerName("selector"),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpcserver.WithDiscovery(r),
		rpcserver.WithEndpoint(goodsserviceName),
		//rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),		// 这是自己封装的 链路追踪
		rpcserver.WithClientEnableTracing(true),
		//rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		panic(err)
	}

	c := goods_pb.NewGoodsClient(conn)
	return c
}
