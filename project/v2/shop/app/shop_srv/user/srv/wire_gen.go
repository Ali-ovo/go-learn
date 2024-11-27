// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package srv

import (
	"shop/app/shop_srv/user/srv/internal/controller/v1"
	"shop/app/shop_srv/user/srv/internal/data/v1/db"
	"shop/app/shop_srv/user/srv/internal/service/v1"
	"shop/gmicro/app"
	"shop/gmicro/pkg/log"
	"shop/pkg/options"
)

// Injectors from wire.go:

func initApp(logOptions *log.Options, serverOptions *options.ServerOptions, registryOptions *options.RegistryOptions, telemetryOptions *options.TelemetryOptions, mySQLOptions *options.MySQLOptions, nacosOptions *options.NacosOptions) (*app.App, error) {
	dataFactory, err := db.GetDBfactoryOr(mySQLOptions)
	if err != nil {
		return nil, err
	}
	serviceFactory := srv.NewService(dataFactory)
	userServer := controller.NewUserServer(serviceFactory)
	nacosDataSource, err := NewNacosDataSource(nacosOptions)
	if err != nil {
		return nil, err
	}
	server, err := NewUserRPCServer(telemetryOptions, serverOptions, userServer, nacosDataSource)
	if err != nil {
		return nil, err
	}
	registrar := NewRegistrar(registryOptions, serverOptions)
	appApp, err := NewUserApp(logOptions, serverOptions, server, registrar)
	if err != nil {
		return nil, err
	}
	return appApp, nil
}