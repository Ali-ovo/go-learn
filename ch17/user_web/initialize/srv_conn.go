package initialize

import (
	"fmt"
	"go-learn/ch17/user_web/global"
	"go-learn/ch17/user_web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() {
	// 从注册中心获取服务器的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		zap.S().Infof("userSrvHost: %s, userSrvPort: %d", userSrvHost, userSrvPort)
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 链接 [用户服务失败]")
		return
	}

	// 链接用户 grpc
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Error("[GetUserList] 链接 [用户服务失败]", "msg", err.Error())
	}

	// 初始化客户端
	userSrvClient := proto.NewUserClient(conn)

	global.UserSrvClient = userSrvClient
}
