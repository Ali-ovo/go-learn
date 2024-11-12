package consul

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"shop/gmicro/pkg/log"
	"strconv"
	"strings"
	"time"

	"shop/gmicro/registry"

	"github.com/hashicorp/consul/api"
)

// Client is consul client config
type Client struct {
	cli    *api.Client
	ctx    context.Context
	cancel context.CancelFunc // 用来 退出程序

	// resolve service entry endpoints 解析 服务入口端点
	resolver ServiceResolver
	// healthcheck time interval in seconds 健康检查 以秒为单位时间间隔  服务端 -> 客户端 你死没死的消息
	healthcheckInterval int
	// heartbeat enable heartbeat 是否开启心跳  客户端 -> 服务端 我还存活着的消息
	heartbeat bool
	// deregisterCriticalServiceAfter time interval in seconds 注销服务请求发起后 以秒为单位时间间隔 的超时时间
	deregisterCriticalServiceAfter int
	// serviceChecks  user custom checks 自定义检查接口
	serviceChecks api.AgentServiceChecks
}

// NewClient creates consul client
func NewClient(cli *api.Client) *Client {
	c := &Client{
		cli:                            cli,
		resolver:                       defaultResolver,
		healthcheckInterval:            10,
		heartbeat:                      true,
		deregisterCriticalServiceAfter: 600,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries)) // 创建 一个列表 存储相关信息
	for _, entry := range entries {
		var version string
		for _, tag := range entry.Service.Tags { // 例如: "Tags": ["version=2.1.1"],
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port))
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}

	return services
}

// ServiceResolver is used to resolve service endpoints
type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance

// Service get services from consul
func (c *Client) Service(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := c.cli.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}
	return c.resolver(ctx, entries), meta.LastIndex, nil
}

// Register register service instance to consul 注册 consul 服务实例
func (c *Client) Register(_ context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	//
	addresses := make(map[string]api.ServiceAddress, len(svc.Endpoints))
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints { // 遍历 终端 server path
		raw, err := url.Parse(endpoint) // 解析 成为 url 格式
		if err != nil {
			return err
		}
		addr := raw.Hostname() // 获得 host
		// 解析成 10进制 Uint16 的类型
		port, _ := strconv.ParseUint(raw.Port(), 10, 16) // 保证 port 的准确性

		checkAddresses = append(checkAddresses, net.JoinHostPort(addr, strconv.FormatUint(port, 10))) // 将修改完成的 host:port 添加
		// 例如 addresses["http"] addresses["grpc"]
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}
	// 是否开启健康检查
	if enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            address,
				Interval:                       fmt.Sprintf("%ds", c.healthcheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
				Timeout:                        "5s",
			})
		}
	}
	// 是否开启心跳检查
	if c.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID: "service:" + svc.ID,
			// 指定这是一个 TTL 检查，必须定期使用 TTL 端点来更新检查状态。如果检查未设置为在指定时间内通过，则检查将设置为失败状态。
			TTL:                            fmt.Sprintf("%ds", c.healthcheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
		})
	}

	// custom checks 自定义检查
	asr.Checks = append(asr.Checks, c.serviceChecks...)

	err := c.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	// 实现 心跳
	if c.heartbeat {
		go func() {
			time.Sleep(time.Second)
			err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			if err != nil {
				log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
			}
			ticker := time.NewTicker(time.Second * time.Duration(c.healthcheckInterval)) // 用于执行定时操作
			defer ticker.Stop()                                                          // 结束定时操作
			for {
				select {
				case <-ticker.C:
					err = c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass") // 发送心跳
					if err != nil {
						log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
					}
				case <-c.ctx.Done(): // 外部想要结束心跳 需要 执行 ctx.done()
					return
				}
			}
		}()
	}
	return nil
}

// Deregister deregister service by service ID 通过 service ID 注销 服务
func (c *Client) Deregister(_ context.Context, serviceID string) error {
	c.cancel()
	return c.cli.Agent().ServiceDeregister(serviceID)
}
