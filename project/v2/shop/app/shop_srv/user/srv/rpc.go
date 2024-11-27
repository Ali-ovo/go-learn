package srv

import (
	"fmt"
	user_pb "shop/api/user/v1"
	"shop/gmicro/core/trace"
	"shop/gmicro/pkg/conn"
	"shop/gmicro/server/rpcserver"
	"shop/pkg/options"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
	"github.com/alibaba/sentinel-golang/pkg/datasource/nacos"
)

func NewNacosDataSource(opts *options.NacosOptions) (*nacos.NacosDataSource, error) {
	switch opts.LogLevel {
	case "debug":
		logging.ResetGlobalLoggerLevel(logging.DebugLevel)
	case "info":
		logging.ResetGlobalLoggerLevel(logging.InfoLevel)
	case "warn":
		logging.ResetGlobalLoggerLevel(logging.WarnLevel)
	case "error":
		logging.ResetGlobalLoggerLevel(logging.ErrorLevel)
	}
	configClient, err := conn.NewNacosClient(&conn.NacosOptions{
		Host:                 opts.Host,
		Port:                 uint64(opts.Port),
		User:                 opts.User,
		Password:             opts.Password,
		TimeOut:              opts.TimeOut,
		NotLoadCacheAtStart:  false,
		UpdateCacheWhenEmpty: false,
		LogDir:               opts.LogDir,
		LogLevel:             opts.LogLevel,
		CacheDir:             opts.CacheDir,
		NamespaceId:          opts.NameSpace,
	})
	if err != nil {
		return nil, err
	}

	// 注册 流控规则
	h := datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser)
	// 创建 NacosDataSource 数据源
	return nacos.NewNacosDataSource(configClient, opts.Group, opts.DataID, h)
}

func NewUserRPCServer(telemetry *options.TelemetryOptions, serverOpts *options.ServerOptions, userServer user_pb.UserServer, dataNacos *nacos.NacosDataSource) (*rpcserver.Server, error) {
	// 初始化 open-telemetry的 exporter
	trace.InitAgent(trace.Options{
		Name:     telemetry.Name,
		Endpoint: telemetry.Endpoint,
		Sampler:  telemetry.Sampler,
		Batcher:  telemetry.Batcher,
	})

	rpcAddr := fmt.Sprintf("%s:%d", serverOpts.Host, serverOpts.Port)

	var opts []rpcserver.ServerOption
	opts = append(
		opts,
		rpcserver.WithAddress(rpcAddr),
		rpcserver.WithServerMetrics(serverOpts.EnableMetrics),
		rpcserver.WithServerEnableTracing(serverOpts.EnableTelemetry),
	)
	if serverOpts.EnableLimit {
		opts = append(opts, rpcserver.WithServerUnaryInterceptor(grpc.NewUnaryServerInterceptor()))
		err := dataNacos.Initialize()
		if err != nil {
			return nil, err
		}
	}

	urpcServer := rpcserver.NewServer(opts...)
	user_pb.RegisterUserServer(urpcServer.Server, userServer)

	return urpcServer, nil
}

//func NewUserRPCServer(cfg *config.Config) (*rpcserver.Server, error) {
//	// 初始化 open-telemetry的 exporter
//	trace.InitAgent(trace.Options{
//		Name:     cfg.Telemetry.Name,
//		Endpoint: cfg.Telemetry.Endpoint,
//		Sampler:  cfg.Telemetry.Sampler,
//		Batcher:  cfg.Telemetry.Batcher,
//	})
//
//	// 有点繁琐, wire, ioc-golang
//	dataFactory, err := db.GetDBfactoryOr(cfg.Mysql)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	srvFactory := srv.NewService(dataFactory)
//	iServer := controller.NewUserServer(srvFactory)
//
//	rpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
//	urpcServer := rpcserver.NewServer(
//		rpcserver.WithAddress(rpcAddr),
//		rpcserver.WithServerMetrics(cfg.Server.EnableMetrics),
//		rpcserver.WithServerEnableTracing(cfg.Server.EnableTelemetry),
//	)
//	user_pb.RegisterUserServer(urpcServer.Server, iServer)
//
//	return urpcServer, nil
//}
