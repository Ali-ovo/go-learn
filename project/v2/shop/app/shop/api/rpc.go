package api

import (
	"fmt"
	goods_pb "shop/api/goods/v1"
	"shop/app/shop/api/config"
	"shop/app/shop/api/internal_srv/controller/v1"
	"shop/app/shop/api/internal_srv/data/v1/db"
	"shop/app/shop/api/internal_srv/data_search/v1/es"
	srv "shop/app/shop/api/internal_srv/service/v2"
	"shop/gmicro/core/trace"
	"shop/gmicro/pkg/log"
	"shop/gmicro/server/rpcserver"
)

func NewGoodsRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	// 数据库的工厂方法
	dataFactory, err := db.GetDBfactoryOr(cfg.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	searchFactory, err := es.GetSearchFactoryOr(cfg.EsOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 业务层的工厂方法
	srvFactory := srv.NewService(dataFactory, searchFactory)

	gServer := controller.NewGoodsServer(srvFactory)

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	gRpcServer := rpcserver.NewServer(
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
	)
	goods_pb.RegisterGoodsServer(gRpcServer.Server, gServer)

	return gRpcServer, nil
}
