package user

import (
	"context"
	user_pb "shop/api/user/v1"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver"
)

const userserviceName = "discovery:///shop-user-srv"

func NewUserServiceClient(r registry.Discovery) user_pb.UserClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		// 设置负载均衡
		rpcserver.WithBanlancerName("selector"),
		rpcserver.WithDiscovery(r),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpcserver.WithEndpoint(userserviceName),
		//rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),		// 这是自己封装的 链路追踪
		rpcserver.WithClientEnableTracing(true),
		//rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)
	if err != nil {
		panic(err)
	}

	c := user_pb.NewUserClient(conn)
	return c
}
