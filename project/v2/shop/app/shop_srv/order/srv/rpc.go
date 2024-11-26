package srv

import (
	"fmt"
	order_pb "shop/api/order/v1"
	"shop/app/shop_srv/order/srv/config"
	"shop/app/shop_srv/order/srv/internal/controller/v1"
	"shop/app/shop_srv/order/srv/internal/data/v1/db"
	"shop/app/shop_srv/order/srv/internal/service/v1/srv"
	"shop/gmicro/core/trace"
	"shop/gmicro/pkg/log"
	"shop/gmicro/server/rpcserver"
)

func NewOrderRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	// 数据库的工厂方法
	dataFactory, err := db.GetDataFactoryOr(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, cfg.Dtm)
	iServer := controller.NewOrderServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	iRpcServer := rpcserver.NewServer(
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
	)
	order_pb.RegisterOrderServer(iRpcServer.Server, iServer)
	return iRpcServer, nil
}
