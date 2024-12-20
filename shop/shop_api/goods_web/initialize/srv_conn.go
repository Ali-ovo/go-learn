package initialize

import (
	"fmt"
	"go-learn/shop/shop_api/goods_web/global"
	"go-learn/shop/shop_api/goods_web/proto"

	"go-learn/shop/shop_api/goods_web/utils/otgrpc"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 链接 [用户服务失败]: %v", err)
	}

	// 初始化客户端
	global.GoodsSrvClient = proto.NewGoodsClient(conn)
}
