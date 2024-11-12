package app

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"shop/gmicro/pkg/log"
	"shop/gmicro/registry"
	"sync"
)

type App struct {
	//logger *log.Logger
	lk       sync.Mutex // struct 类型 说明是实例化后的  可以直接用 不需要再次声明
	opts     options
	instance *registry.ServiceInstance // registry 参数
}

func New(opts ...Option) *App {
	//o := options{
	//	sigs:             []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT},
	//	registrarTimeout: 10 * time.Second, // 注册服务 超时时间
	//	stopTimeout:      10 * time.Second, // 注销服务 超时时间
	//}
	//if id, err := uuid.NewUUID(); err == nil {
	//	o.id = id.String()
	//}
	o := DefaultOptions()

	for _, opt := range opts { // 执行附加参数
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

// Run 启动整个微服务
func (a *App) Run() error {
	// 注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	// 这个变量可能被其他的 goroutine 访问到
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	if a.opts.rpcServer != nil {
		go func() {
			err := a.opts.rpcServer.Start()
			if err != nil {
				panic(err)
			}
		}()
	}

	// 注册 服务 到 consul 等...
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err := a.opts.registrar.Register(rctx, a.instance)
		if err != nil {
			log.Errorf("register service error: %s", err)
			return err
		}
	}

	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	<-c
	return nil
}

// Stop 停止服务
func (a *App) Stop() error {
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()
	if a.opts.registrar != nil && instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err := a.opts.registrar.Deregister(rctx, a.instance)
		if err != nil {
			log.Errorf("deregister service error: %s", err)
			return err
		}
	}
	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}

	// 从rpcserver, restserver 去主动获取这些信息
	if a.opts.rpcServer != nil {
		// u := a.opts.rpcServer.Endpoint()
		u := &url.URL{
			Scheme: "grpc",
			Host:   a.opts.rpcServer.Address(),
		}
		endpoints = append(endpoints, u.String())
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
	}, nil
}
