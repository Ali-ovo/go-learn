package rpcserver

import (
	"shop/gmicro/registry"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
	"shop/gmicro/server/rpcserver/selector"
)

const (
	balancerName = "selector"
)

var (
	_ base.PickerBuilder = &balancerBuilder{}
	_ balancer.Picker    = &balancerPicker{}
)

func InitBuilder() {
	// 使用 grpc 提前封装好的 baseBuilder
	b := base.NewBalancerBuilder(
		balancerName, // 负载均衡器名称
		&balancerBuilder{ // 实现 Build方法的结构体
			builder: selector.GlobalSelector(),
		},
		base.Config{HealthCheck: true},
	)
	balancer.Register(b) // 全局注册 负载均衡器
}

type balancerBuilder struct {
	builder selector.Builder
}

// Build creates a grpc Picker.
// info base.PickerBuildInfo 传递是可以直接用 空闲的子链接
func (b *balancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		// Block the RPC until a new picker is available via UpdateState().
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]selector.Node, 0, len(info.ReadySCs)) // 自己进行组装
	for conn, info := range info.ReadySCs {
		// info.Address.Attributes 是由我们传递的数据 有更多的数据 进行负载均衡  服务发现 设置了值
		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*registry.ServiceInstance)
		nodes = append(nodes, &grpcNode{
			// 将 地址 和 我们传递的信息包装到一起(最重要的是 可以在携带的元数据中 传递初始权重)
			Node:    selector.NewNode("grpc", info.Address.Addr, ins),
			subConn: conn, // 空闲的子连接
		})
	}
	p := &balancerPicker{
		selector: b.builder.Build(),
	}
	p.selector.Apply(nodes)
	return p
}

// balancerPicker is a grpc picker.
type balancerPicker struct {
	selector selector.Selector
}

// Pick instances.
func (p *balancerPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	n, done, err := p.selector.Select(info.Ctx)
	if err != nil {
		return balancer.PickResult{}, err
	}

	return balancer.PickResult{
		SubConn: n.(*grpcNode).subConn,
		// Done func(DoneInfo) 格式 进行回调
		Done: func(di balancer.DoneInfo) {
			done(info.Ctx, selector.DoneInfo{
				Err:           di.Err,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ReplyMD:       Trailer(di.Trailer),
			})
		},
	}, nil
}

// Trailer is a grpc trailder MD.
type Trailer metadata.MD

// Get get a grpc trailer value.
func (t Trailer) Get(k string) string {
	v := metadata.MD(t).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

type grpcNode struct {
	selector.Node
	subConn balancer.SubConn
}
