package initialize

import (
	"fmt"
	"go-learn/shop/shop_srvs/order_srv/global"
	"go-learn/shop/shop_srvs/order_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	// 初始化其他服务链接 client
	consulInfo := global.ServerConfig.ConsulInfo

	goodsConn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatalf("初始化链接 [商品服务失败]: %v", err)
	}

	inventoryConn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatalf("初始化链接 [库存服务失败]: %v", err)
	}

	// 初始化客户端
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}
