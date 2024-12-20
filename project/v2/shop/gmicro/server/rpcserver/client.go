package rpcserver

import (
	"context"
	"crypto/tls"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver/clientinterceptors"
	"shop/gmicro/server/rpcserver/resolver/discovery"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	tlsConf  *tls.Config
	timeout  time.Duration
	// discovery 接口
	discovery    registry.Discovery             // 服务发现
	unaryIntes   []grpc.UnaryClientInterceptor  // 一元拦截器
	streamIntes  []grpc.StreamClientInterceptor // 流式拦截器
	rpcOpts      []grpc.DialOption
	balancerName string
	//log          log.LogHelper
	enableTracing bool
	enableMetrics bool
}

// WithTLSConfig with tls config.
func WithTLSConfig(c *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = c
	}
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

// WithClientEnableTracing 设置是否开启链路追踪
func WithClientEnableTracing(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = enable
	}
}

// WithClientEnableMetrics 设置是否开启 普罗米修斯 监控
func WithClientEnableMetrics(enable bool) ClientOption {
	return func(o *clientOptions) {
		o.enableMetrics = enable
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
		timeout:       2000 * time.Millisecond,
		balancerName:  "round_robin", // 默认 grpc 的轮询
		enableTracing: false,
	}

	for _, o := range opts {
		o(&options)
	}

	// 客户端默认拦截器
	ints := []grpc.UnaryClientInterceptor{
		clientinterceptors.TimeoutInterceptor(options.timeout),
	}

	// 开启 链路追踪
	if options.enableTracing {
		ints = append(ints, otelgrpc.UnaryClientInterceptor())
	}

	// 开启 普罗米修斯
	if options.enableMetrics {
		ints = append(ints, clientinterceptors.UnaryPrometheusInterceptor)
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

	// 服务发现的选项
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers( // 添加解析器	参数需要 resolver.Builder
			// 这里创建了一个 结构体 resolver.Builder 然后 grpc 会调用 这个结构的 build 方法 resolver.Builder.Build 实现 服务发现的功能
			discovery.NewBuilder(
				options.discovery,
				discovery.WithInsecure(insecure),
			),
		))
	}

	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	}
	if options.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.rpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.rpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
