package srv

import (
	"context"
	"shop/app/shop_srv/inventory/srv/config"
	gapp "shop/gmicro/app"
	"shop/gmicro/pkg/app"
	"shop/gmicro/pkg/log"
	"shop/gmicro/pkg/storage"
	"shop/gmicro/registry"
	"shop/gmicro/registry/consul"
	"shop/pkg/options"
	"sync"

	"github.com/hashicorp/consul/api"
)

// NewApp 会读取相关配置 (用来集成 pkg/app 用来完成 外部参数校验和映射)
func NewApp(basename string) *app.App {
	cfg := config.NewConfig()

	return app.NewApp(
		"inventory-srv",
		basename,
		app.WithOptions(cfg), // 初始 log server 配置
		app.WithRunFunc(run(cfg)),
		// go run .\main.go --server.port=8081 --server.host=192.168.16.151 --consul.address=192.168.16.105:8500
		//app.WithNoConfig(), // 设置不读取配置文件
	)
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		userApp, err := NewInventoryApp(cfg)
		if err != nil {
			return err
		}

		// 启动 RPC 服务
		if err := userApp.Run(); err != nil {
			log.Errorf("run inventory app error: %s", err)
			return err
		}
		return nil
	}
}

func NewInventoryApp(cfg *config.Config) (*gapp.App, error) {
	// 初始化 log
	log.Init(cfg.Log)
	defer log.Flush()

	// 服务注册
	register := NewRegistrar(cfg.Registry, cfg.Server)

	// 连接redis
	redisConfig := &storage.Config{
		Host:                  cfg.Redis.Host,
		Port:                  cfg.Redis.Port,
		Addrs:                 cfg.Redis.Addrs,
		MasterName:            cfg.Redis.MasterName,
		Username:              cfg.Redis.Username,
		Password:              cfg.Redis.Password,
		Database:              cfg.Redis.Database,
		MaxIdle:               cfg.Redis.MaxIdle,
		MaxActive:             cfg.Redis.MaxActive,
		Timeout:               cfg.Redis.TimeOut,
		EnableCluster:         cfg.Redis.EnableCluster,
		UseSSL:                cfg.Redis.UseSSL,
		SSLInsecureSkipVerify: cfg.Redis.SSLInsecureSkipVerify,
		EnableTracing:         cfg.Redis.EnableTracing,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go storage.ConnectToRedis(context.Background(), &wg, redisConfig)
	wg.Wait()

	// 生成 rpc 服务
	rpcServer, err := NewInventoryRPCServer(cfg)
	if err != nil {
		return nil, err
	}
	// 运行 RPC 服务
	return gapp.New(
		gapp.WithName(cfg.Server.Name),
		gapp.WithServer(rpcServer),
		gapp.WithRegistrar(register),
	), nil
}

func NewRegistrar(registry *options.RegistryOptions, server *options.ServerOptions) registry.Registrar {
	c := api.DefaultConfig()
	c.Address = registry.Address
	c.Scheme = registry.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	return consul.New(cli, consul.WithHealthCheck(server.EnableHealthCheck))
}
