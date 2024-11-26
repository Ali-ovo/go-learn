//go:build wireinject
// +build wireinject

package srv

import (
	"shop/app/shop_srv/user/srv/config"
	gapp "shop/gmicro/app"

	"github.com/google/wire"
)

func initApp(*config.Config) (*gapp.App, error) {
	wire.Build(NewUserApp)
	return &gapp.App{}, nil
}
