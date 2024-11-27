package srv

import (
	"fmt"
	user_pb "shop/api/user/v1"
	"shop/gmicro/core/trace"
	"shop/gmicro/server/rpcserver"
	"shop/pkg/options"

	"github.com/alibaba/sentinel-golang/pkg/adapters/grpc"
)

func NewUserRPCServer(telemetry *options.TelemetryOptions, serverOpts *options.ServerOptions, userServer user_pb.UserServer) (*rpcserver.Server, error) {
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
