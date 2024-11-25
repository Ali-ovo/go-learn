package api

import (
	"fmt"
	"shop/app/shop/api/config"
	"shop/gmicro/core/trace"
	"shop/gmicro/server/restserver"
)

func NewAPIHTTPServer(cfg *config.Config) (*restserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     cfg.Telemetry.Name,
		Endpoint: cfg.Telemetry.Endpoint,
		Sampler:  cfg.Telemetry.Sampler,
		Batcher:  cfg.Telemetry.Batcher,
	})

	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HttpPort)
	aRestServer := restserver.NewServer(
		restserver.WithMode(cfg.Server.HttpMode),                       // 设置 gin 模式
		restserver.WithAddress(rpcAddr),                                // 设置地址
		restserver.WithMiddlewares(cfg.Server.Middlewares),             // 设置中间件
		restserver.WithClientEnableTracing(cfg.Server.EnableTelemetry), // 开启 open-telemetry
		restserver.WithEnableMetrics(cfg.Server.EnableMetrics),         // 是否开启 普罗米修斯监控 默认开启
		restserver.WithEnableProfiling(cfg.Server.EnableProfiling),     // 是否开启 pprof接口 不设置默认开启
	)

	// 配置好路由
	initRouter(aRestServer, cfg)
	return aRestServer, nil
}
