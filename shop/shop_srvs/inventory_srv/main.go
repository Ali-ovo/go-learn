package main

import (
	"flag"
	"fmt"
	"go-learn/shop/shop_srvs/inventory_srv/global"
	"go-learn/shop/shop_srvs/inventory_srv/handler"
	"go-learn/shop/shop_srvs/inventory_srv/initialize"
	"go-learn/shop/shop_srvs/inventory_srv/proto"
	"go-learn/shop/shop_srvs/inventory_srv/utils"
	"go-learn/shop/shop_srvs/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip 地址")
	Port := flag.Int("port", 0, "端口号")

	// 初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("ip: ", *IP, "  port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	if err != nil {
		panic(err)
	}

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	serviceId := uuid.NewV4().String()
	err = registerClient.Register(
		global.ServerConfig.Host,
		*Port,
		global.ServerConfig.Name,
		global.ServerConfig.Tags,
		serviceId,
	)

	if err != nil {
		zap.S().Panic("注册服务失败", err.Error())
	}

	zap.S().Infof("启动服务, 端口: %d", *Port)

	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	// 监听归还topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.189.128:9876"}),
		consumer.WithGroupName("shop_inventory"),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("订阅失败")
	}

	_ = c.Start()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("服务注销失败")
		return
	}
	zap.S().Info("注销成功")
}
