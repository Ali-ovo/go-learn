package srv

import (
	"fmt"
	inventory_pb "shop/api/inventory/v1"
	"shop/app/shop_srv/inventory/srv/config"
	"shop/app/shop_srv/inventory/srv/internal/controller/v1"
	"shop/app/shop_srv/inventory/srv/internal/data/v1/db"
	srv "shop/app/shop_srv/inventory/srv/internal/service/v2"
	"shop/gmicro/core/trace"
	"shop/gmicro/pkg/log"
	"shop/gmicro/server/rpcserver"
)

func NewInventoryRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	// 数据库的工厂方法
	dataFactory, err := db.GetDBfactoryOr(cfg.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory)
	iServer := controller.NewInventoryServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	iRpcServer := rpcserver.NewServer(
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
	)
	inventory_pb.RegisterInventoryServer(iRpcServer.Server, iServer)
	return iRpcServer, nil
}
