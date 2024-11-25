package rpcserver

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"shop/gmicro/api/metadata"
	"shop/gmicro/pkg/common/endpoint"
	"shop/gmicro/pkg/host"
	"shop/gmicro/pkg/log"
	srvintc "shop/gmicro/server/rpcserver/serverinterceptors"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server
	address     string
	lis         net.Listener
	tlsConf     *tls.Config
	endpoint    *url.URL
	err         error
	unaryIntes  []grpc.UnaryServerInterceptor  // 一元拦截器
	streamIntes []grpc.StreamServerInterceptor // 流式拦截器
	grpcOpts    []grpc.ServerOption
	timeout     time.Duration
	//health        *health.Server
	metadata      *metadata.Server
	enableTracing bool
	enableMetrics bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0", // 在没有设置 address 自己获取 ip 和 端口号
		//health:  health.NewServer(),
		//timeout: 1 * time.Second,
		enableTracing: true,
	}

	for _, opt := range opts {
		opt(srv)
	}

	// TODO 希望用户不设置拦截器的情况下, 会自动默认加上一些必须的拦截器, crash, tracingtry
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryRecoverInterceptor, // 一元拦截器 异常处理(而不是一层层抛出 然后停止程序)
	}
	// 开启链路追踪
	if srv.enableTracing {
		unaryInts = append(unaryInts, otelgrpc.UnaryServerInterceptor())
	}
	// 开启 普罗米修斯监控
	if srv.enableMetrics {
		unaryInts = append(unaryInts, srvintc.UnaryPrometheusInterceptor)
	}

	if srv.timeout > 0 {
		unaryInts = append(unaryInts, srvintc.UnaryTimeoutInterceptor(srv.timeout))
	}

	if len(srv.unaryIntes) > 0 {
		unaryInts = append(unaryInts, srv.unaryIntes...)
	}
	streamIntes := []grpc.StreamServerInterceptor{
		srvintc.StreamRecoverInterceptor, // 流式拦截器 错误 不退出程序 继续运行
	}
	if len(srv.streamIntes) > 0 {
		streamIntes = append(streamIntes, srv.streamIntes...)
	}

	// 把传入的拦截器 转换成 grpc 的 ServerOption
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(unaryInts...), grpc.ChainStreamInterceptor(streamIntes...)}

	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}

	// 把用户传入的grpc.ServerOption 放在一起
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)
	// register metadata Server
	srv.metadata = metadata.NewServer(srv.Server)
	// analysis 解析 address
	if err := srv.listenAndEndpoint(); err != nil {
		panic(err)
	}

	//// register 注册 grpc health
	//grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)// 使用grpc 健康检查的时候用  无grpc 的健康检查没有用 (本项目没有使用 直接使用 TCP连接)

	// 可以支持 用户直接通过 grpc 的一个接口查看当前支持的所有的 rpc 服务
	metadata.RegisterMetadataServer(srv.Server, srv.metadata) // 关键函数 是gRPC服务器端用于注册服务的方法
	reflection.Register(srv.Server)                           // TODO ??? 不知道有什么用 需要配合 grpc.client 也需要写对于代码
	return srv
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func WithLis(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func WithServerUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *Server) {
		o.unaryIntes = in
	}
}

func WithStreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *Server) {
		o.streamIntes = in
	}
}

func WithServerOption(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

//func WithHealthServer(health *health.Server) ServerOption {
//	return func(s *Server) {
//		s.health = health
//	}
//}

func WithServerTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// WithServerEnableTracing 设置是否开启链路追踪
func WithServerEnableTracing(enableTracing bool) ServerOption {
	return func(s *Server) {
		s.enableTracing = enableTracing
	}
}

// WithServerMetrics 设置是否开启 普罗米修斯 监控
func WithServerMetrics(metric bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = metric
	}
}

// Start 启动 grpc 的服务
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	//s.health.Resume() // 这里可以不写  只不过 我这里显示声明一下(健康检查服务知道该服务已经恢复正常工作)
	return s.Server.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	//s.health.Shutdown() // 设置 服务的状态为 not_serving, 阻止 接收新的请求 过来
	s.GracefulStop() // grpc 优雅退出
	log.Infof("[gRPC] server stopping")
	return nil
}

func (s *Server) Address() string {
	return s.address
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// 完成 ip 和 端口的提取
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return nil
}
