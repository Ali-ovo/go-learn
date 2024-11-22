package admin

import (
	"shop/app/shop_api/admin/config"

	"shop/gmicro/server/restserver"
)

func NewUserHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	urestServer := restserver.NewServer(
		restserver.WithPort(cfg.Server.HttpPort),
		restserver.WithMiddlewares(cfg.Server.Middlewares),
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics),
	)

	// 配置好路由
	initRouter(urestServer)
	return urestServer, nil
}