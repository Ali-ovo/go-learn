package registry

import "context"

// Registrar Service registration interface 服务注册接口
type Registrar interface {
	// Register 注册
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister 注销
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discovery Service Discovery Interface
type Discovery interface {
	// GetService 获取服务实例
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch 创建服务监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

type Watcher interface {
	// Next 获取服务实例, next在下面的情况下会返回服务
	//1. 第一次监听时，如果服务实例列表不为空，则返回服务实例列表
	//2. 如果服务实例发生变化，则返回服务实例列表
	//3. 如果上面两种情况都不满足，则会阻塞到context deadline或者cancel
	Next() ([]*ServiceInstance, error)
	// Stop actively give up listening 主动放弃监听
	Stop() error
}

type ServiceInstance struct {
	// Service ID registered with the registry center. 注册到注册中心的服务id
	ID string `json:"id"`
	// Service name
	Name string `json:"name"`
	// Service version
	Version string `json:"version"`
	// Service metadata 服务元数据
	Metadata map[string]string `json:"metadata"`
	// Service path
	// http://127.0.0.1:8000
	// grpc://127.0.0.1:9000
	Endpoints []string `json:"endpoints"`
}
