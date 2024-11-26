package consul

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"shop/gmicro/registry"

	"github.com/hashicorp/consul/api"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is consul registry option.
type Option func(*Registry)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// WithHeartbeat enable or disable heartbeat
func WithHeartbeat(enable bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.heartbeat = enable
		}
	}
}

// WithServiceResolver with endpoint function option.
func WithServiceResolver(fn ServiceResolver) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.resolver = fn
		}
	}
}

// WithHealthCheckInterval with healthcheck interval in seconds.
func WithHealthCheckInterval(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.healthcheckInterval = interval
		}
	}
}

// WithDeregisterCriticalServiceAfter with deregister-critical-service-after in seconds.
func WithDeregisterCriticalServiceAfter(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.deregisterCriticalServiceAfter = interval
		}
	}
}

// WithServiceCheck with service checks
func WithServiceCheck(checks ...*api.AgentServiceCheck) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.serviceChecks = checks
		}
	}
}

// Config is consul registry config
type Config struct {
	*api.Config
}

// Registry is consul registry
type Registry struct {
	cli               *Client
	enableHealthCheck bool // 是否开启健康检查
	registry          map[string]*serviceSet
	lock              sync.RWMutex // 读写锁
}

// New creates consul registry
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		cli:               NewClient(apiClient),
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Register register service
func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

// Deregister deregister service
func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService return service by name
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	set := r.registry[name]

	getRemote := func() []*registry.ServiceInstance {
		services, _, err := r.cli.Service(ctx, name, 0, true)
		if err == nil && len(services) > 0 {
			return services
		}
		return nil
	}

	if set == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if ss == nil {
		if s := getRemote(); len(s) > 0 {
			return s, nil
		}
		return nil, fmt.Errorf("service %s not found in registry", name)
	}
	return ss, nil
}

// ListServices return service list.
// ListServices 返回 srv 列表
func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		var services []*registry.ServiceInstance
		ss, _ := set.services.Load().([]*registry.ServiceInstance)
		if ss == nil {
			continue
		}
		services = append(services, ss...)
		allServices[name] = services
	}
	return
}

// Watch resolve service by name 通过 名称 解析服务
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	// 锁此函数 防止 多个 协程 导致 数据不一致
	r.lock.Lock()
	defer r.lock.Unlock()
	set, ok := r.registry[name] // 如果 registry 中无数据 创建一个 set 第一次访问 为空
	if !ok {
		set = &serviceSet{
			watcher:     make(map[*watcher]struct{}), // 创建 观察者
			services:    &atomic.Value{},             // 原子性操作
			serviceName: name,
		}
		r.registry[name] = set
	}

	// 初始化 watcher 观察者
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{} // 将 实例 w 写入到 set 中
	set.lock.Unlock()
	ss, _ := set.services.Load().([]*registry.ServiceInstance) // 原子性读	断言失败为 nil
	if len(ss) > 0 {
		// If the service has a value, it needs to be pushed to the watcher,
		// otherwise the initial data may be blocked forever during the watch.
		// 如果服务具有值，则需要将其推送到观察者，否则在观察期间初始数据可能会被永久阻塞
		// 说明 set.services 里已经存在值了  并且有个 协程在 不停的运行
		w.event <- struct{}{} // 从 channel 中写入值 如果 channel 有值 就会 hook 住
	}

	if !ok { // 第一次执行
		err := r.resolve(set)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (r *Registry) resolve(ss *serviceSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	services, idx, err := r.cli.Service(ctx, ss.serviceName, 0, true) // ss.serviceName 服务注册名称
	cancel()
	if err != nil {
		return err
	} else if len(services) > 0 { // 获取到服务了
		ss.broadcast(services) // 原子性写操作 存入 consul 返回回来的相关信息	[]*registry.ServiceInstance
	}
	go func() {
		ticker := time.NewTicker(time.Second) // 执行定时任务
		defer ticker.Stop()
		for {
			<-ticker.C
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*120) // 长轮询的思想
			tmpService, tmpIdx, err := r.cli.Service(ctx, ss.serviceName, idx, true)
			cancel()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if len(tmpService) != 0 && tmpIdx != idx { // 获取到修改之后的服务了
				services = tmpService  // 健康的服务列表
				ss.broadcast(services) // 进行原子性的存储
			}
			idx = tmpIdx
		}
	}()

	return nil
}
