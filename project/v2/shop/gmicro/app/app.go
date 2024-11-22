package app

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"shop/gmicro/pkg/log"
	"shop/gmicro/registry"
	gs "shop/gmicro/server"
	"sync"

	"golang.org/x/sync/errgroup"
)

type App struct {
	//logger *log.Logger
	lk       sync.Mutex // struct 类型 说明是实例化后的  可以直接用 不需要再次声明
	opts     options
	instance *registry.ServiceInstance // registry 参数
	cancel   context.CancelFunc
}

func New(opts ...Option) *App {
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

	//现在启动了两个server，一个是restserver，一个是rpcserver
	/*
		这两个server是否必须同时启动成功？
		如果有一个启动失败，那么我们就要停止另外一个server
		如果启动了多个， 如果其中一个启动失败，其他的应该被取消
			如果剩余的server的状态：
				1. 还没有开始调用start
					stop
				2. start进行中
					调用进行中的cancel
				3. start已经完成
					调用stop
		如果我们的服务启动了然后这个时候用户立马进行了访问
	*/

	var servers []gs.Server
	if a.opts.restServer != nil {
		servers = append(servers, a.opts.restServer)
	}

	if a.opts.rpcServer != nil {
		servers = append(servers, a.opts.rpcServer)
	}

	// 启动 resetserver
	parentCtx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(parentCtx)
	wg := sync.WaitGroup{}
	for _, srv := range servers {
		// 这样 协程中调用的 srv 就会引用函数内部的变量  不会因为 srv 的改变而改变
		// 不做此操作 有可能发生 下面协程 srv.Start 中 srv 启动的是 其他的 微服务 而没启动 本身的微服务
		srv := srv

		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			// 不可能无休止的等待 stop 信号
			sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})

		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start rest server")
			return srv.Start(ctx)
		})
	}

	wg.Wait() // 上面 api 和 grpc 服务都正常开启后 才能继续运行 否则会 hold 住

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
	eg.Go(func() error {
		select {
		case <-ctx.Done(): // 主动关闭
			return ctx.Err()
		case <-c: // 获取退出信号
			return a.Stop() // 执行退出
		}
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

/*
	http basic 认证
	cache
	cache 1. redis 2.memcache 3. local cache
	jwt
*/

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
	if a.cancel != nil {
		a.cancel()
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
