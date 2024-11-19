package srv

import (
	"fmt"
	upb "shop/api/user/v1"
	"shop/app/user/srv/config"
	"shop/app/user/srv/controller/user"
	"shop/app/user/srv/data/v1/db"
	"shop/gmicro/core/trace"
	"shop/gmicro/server/rpcserver"

	srvv1 "shop/app/user/srv/service/v1"
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
	data := db.NewUsers() // 只操作数据库
	srv := srvv1.NewUserService(data)
	userver := user.NewUserServer(srv)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	urpcServer := rpcserver.NewServer(rpcserver.WithAddress(rpcAddr))
	upb.RegisterUserServer(urpcServer.Server, userver)

	//r := gin.Default()
	//upb.RegisterUserServerHTTPServer(userver, r)
	//r.Run(":8082")
	return urpcServer, nil
}
