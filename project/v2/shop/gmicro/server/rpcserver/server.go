package rpcserver

import (
	"crypto/tls"
	"go-learn/project/v2/shop/pkg/common/endpoint"
	"go-learn/project/v2/shop/pkg/host"
	"net"
	"net/url"
	"time"

	"go-learn/project/v2/shop/gmicro/api/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	gprc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server

	address       string
	lis           net.Listener
	tlsConf       *tls.Config
	endpoint      *url.URL
	err           error
	unaryIntes    []grpc.UnaryServerInterceptor  // 一元拦截器
	streamIntes   []grpc.StreamServerInterceptor // 流式拦截器
	grpcOpts      []grpc.ServerOption
	timeout       time.Duration
	health        *health.Server
	metadata      *metadata.Server
	enableTracing bool
	enableMetrics bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0", // 在没有设置 address 自己获取 ip 和 端口号
		health:  health.NewServer(),
		//timeout: 1 * time.Second,
		enableTracing: true,
	}

	for _, o := range opts {
		o(srv)
	}

	// TODO 默认拦截器

	// 把传入的拦截器转换成 grpc ServerOptions
	grpcOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(srv.unaryIntes...)}

	// 把用户自己传入的 grpc ServerOptions 也加入到 grpcOpts
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	// 创建 grpc 服务
	srv.Server = grpc.NewServer(grpcOpts...)
	srv.metadata = metadata.NewServer(srv.Server)

	// 解析 address
	if err := srv.listenAndEndpoint(); err != nil {
		return nil
	}

	gprc_health_v1.RegisterHealthServer(srv.Server, srv.health)

	metadata.RegisterMetadataServer(srv.Server, srv.metadata)

	reflection.Register(srv.Server)
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
