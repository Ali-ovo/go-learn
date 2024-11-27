//go:build wireinject
// +build wireinject

package srv

import (
	"shop/app/shop_srv/user/srv/internal/controller/v1"
	"shop/app/shop_srv/user/srv/internal/data/v1/db"
	srv "shop/app/shop_srv/user/srv/internal/service/v1"
	gapp "shop/gmicro/app"
	"shop/gmicro/pkg/log"
	"shop/pkg/options"

	"github.com/google/wire"
)

func initApp(*log.Options, *options.ServerOptions, *options.RegistryOptions, *options.TelemetryOptions, *options.MySQLOptions) (*gapp.App, error) {
	// 因为 底层是使用 三层代码架构, wire.NewSet 这个方法了 因为三层代码架构 也可以做到帮我合成一个接口 (个人觉得更适用与微服务)
	wire.Build(ProviderSet, controller.NewUserServer, srv.NewService, db.GetDBfactoryOr)
	return &gapp.App{}, nil
}
