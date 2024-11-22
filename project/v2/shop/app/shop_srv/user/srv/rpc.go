package srv

import (
	"fmt"
	"log"
	upb "shop/api/user/v1"

	"shop/gmicro/core/trace"
	"shop/gmicro/server/rpcserver"

	"shop/app/shop_srv/user/srv/config"
	"shop/app/shop_srv/user/srv/controller/user"
	"shop/app/shop_srv/user/srv/data/v1/db"
	"shop/app/shop_srv/user/srv/service/v1"
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
	gormDB, err := db.GetDBfactoryOr(cfg.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := db.NewUsers(gormDB)
	srv := service.NewUserService(data)
	userver := user.NewUserServer(srv)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	urpcServer := rpcserver.NewServer(
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
	)

	upb.RegisterUserServer(urpcServer.Server, userver)

	return urpcServer, nil
}
