package main

import (
	"flag"
	"fmt"
	"go-learn/ch17/user_srv/global"
	"go-learn/ch17/user_srv/handler"
	"go-learn/ch17/user_srv/initialize"
	"go-learn/ch17/user_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip 地址")
	Port := flag.Int("port", 50051, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("ip:", *IP, "port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	// 检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", "192.168.6.124", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"user", "srv", "ali"}
	registration.Address = "192.168.6.124"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		panic(err)
	}

	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
