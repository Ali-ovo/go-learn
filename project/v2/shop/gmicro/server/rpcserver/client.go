package rpcserver

import (
	"context"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	"shop/gmicro/server/rpcserver/resolver/discovery"
	"time"

	"google.golang.org/grpc"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	timeout  time.Duration
	// discovery 接口
	discovery    registry.Discovery             // 服务发现
	unaryIntes   []grpc.UnaryClientInterceptor  // 一元拦截器
	streamIntes  []grpc.StreamClientInterceptor // 流式拦截器
	rpcOpts      []grpc.DialOption
	balancerName string
	//logger        log.LogHelper
}

// WithEndpoint 设置地址
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithClientTimeout 设置超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithDiscovery 设置服务发现
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithClientUnaryInterceptor 设置 unary 拦截器
func WithClientUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryIntes = in
	}
}

// WithClientStreamInterceptor 设置 stream 拦截器
func WithClientStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamIntes = in
	}
}

// WithDialOption 设置 grpc 的 dial 选项
func WithDialOption(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.rpcOpts = opts
	}
}

// WithBanlancerName 设置负载均衡器
func WithBanlancerName(name string) ClientOption {
	return func(o *clientOptions) {
		o.balancerName = name
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:      2000 * time.Millisecond,
		balancerName: "round_robin",
	}

	for _, o := range opts {
		o(&options)
	}

	// 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
	}
	sints := []grpc.StreamClientInterceptor{}

	if len(options.unaryIntes) > 0 {
		ints = append(ints, options.unaryIntes...)
	}
	if len(options.streamIntes) > 0 {
		sints = append(sints, options.streamIntes...)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "` + options.balancerName + `"}`),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(sints...),
	}

	// TODO 服务发现的选项
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers( // 添加解析器	参数需要 resolver.Builder
			discovery.NewBuilder(
				options.discovery,
				discovery.WithInsecure(insecure),
			),
		))
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}

	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
