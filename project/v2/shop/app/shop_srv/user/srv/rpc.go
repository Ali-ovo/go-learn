package srv

import (
	"fmt"
	upb "shop/api/user/v1"
	"shop/app/shop_srv/user/srv/config"
	"shop/app/shop_srv/user/srv/internal/controller/v1"
	"shop/app/shop_srv/user/srv/internal/data/v1/db"
	srv "shop/app/shop_srv/user/srv/internal/service/v1"
	"shop/gmicro/core/trace"
	"shop/gmicro/pkg/log"
	"shop/gmicro/server/rpcserver"
)

func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	// 有点繁琐, wire, ioc-golang
	dataFactory, err := db.GetDBfactoryOr(cfg.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	srvFactory := srv.NewService(dataFactory)
	iServer := controller.NewUserServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
	)
	upb.RegisterUserServer(urpcServer.Server, iServer)

	return urpcServer, nil
}
