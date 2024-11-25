package admin

import (
	"fmt"
	"shop/app/shop_api/admin/config"

	"shop/gmicro/server/restserver"
)

func NewUserHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)
	urestServer := restserver.NewServer(
		restserver.WithAddress(rpcAddr),
		restserver.WithMiddlewares(cfg.Server.Middlewares),
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics),
	)

	// 配置好路由
	initRouter(urestServer)
	return urestServer, nil
}
