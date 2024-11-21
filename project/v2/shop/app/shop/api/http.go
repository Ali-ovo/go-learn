package admin

import (
	"shop/app/shop/api/config"
	"shop/gmicro/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	aRestServer := restserver.NewServer(
		restserver.WithPort(cfg.Server.HttpPort),               // 设置端口号
		restserver.WithMiddlewares(cfg.Server.Middlewares),     // 设置中间件
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics), // 开启 普罗米修斯监控
	)

	// 配置好路由
	initRouter(aRestServer, cfg)
	return aRestServer, nil
}
