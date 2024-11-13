package admin

import (
	"shop/app/user/srv/config"
	"shop/gmicro/server/restserver"
)

func NewUserHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	urestServer := restserver.NewServer(restserver.WithPort(cfg.Server.HttpPort))

	// 配置好路由
	initRouter(urestServer)
	return urestServer, nil
}
