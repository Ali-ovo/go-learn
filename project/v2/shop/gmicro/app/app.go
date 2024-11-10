package app

import (
	"context"
	registry "go-learn/project/v2/shop/gmicro/registery"
	"go-learn/project/v2/shop/pkg/log"
	"os"
	"os/signal"
	"sync"
)

type App struct {
	opts options

	lk sync.Mutex

	instance *registry.ServiceInstance
}

func New(opts ...Option) *App {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return &App{
		opts: o,
	}
}

// 启动整个服务
func (a *App) Run() error {
	// 注册信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err := a.opts.registrar.Register(rctx, instance)

		if err != nil {
			log.Errorf("register service error: %v", err)
			return err
		}
	}

	// 注册服务

	// 监听退出型号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	<-c
	return nil
}

// 停止服务
func (a *App) Stop() error {
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()

	if a.opts.registrar != nil && a.instance != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		if err := a.opts.registrar.Deregister(rctx, instance); err != nil {
			log.Errorf("deregister service error: %v", err)
			return err
		}
	}
	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {

	endPoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endPoints = append(endPoints, e.String())
	}

	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endPoints,
	}, nil
}
