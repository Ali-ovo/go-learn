package rpcserver

import (
	"net"
	"net/url"
	"shop/gmicro/api/metadata"
	"shop/gmicro/pkg/host"
	"shop/gmicro/pkg/log"
	srvintc "shop/gmicro/server/rpcserver/serverinterceptors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server

	address     string
	unaryIntes  []grpc.UnaryServerInterceptor  // 一元拦截器
	streamIntes []grpc.StreamServerInterceptor // 流式拦截器
	grpcOpts    []grpc.ServerOption
	lis         net.Listener
	timeout     time.Duration

	health       *health.Server
	customHealth bool
	metadata     *metadata.Server
	endpoint     *url.URL
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address: ":0", // 在没有设置 address 自己获取 ip 和 端口号
		health:  health.NewServer(),
		//timeout: 1 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	// TODO 希望用户不设置拦截器的情况下, 会自动默认加上一些必须的拦截器, crash, tracingtry
	unaryInts := []grpc.UnaryServerInterceptor{
		srvintc.UnaryRecoverInterceptor, // 一元拦截器 异常处理(而不是一层层抛出 停止程序)
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

	// register 注册 grpc health
	if !srv.customHealth {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	}
	// 可以支持 用户直接通过 grpc 的一个接口查看当前支持的所有的 rpc 服务
	metadata.RegisterMetadataServer(srv.Server, srv.metadata) // 关键函数 是 gRPC 服务器端用于注册服务的方法
	reflection.Register(srv.Server)                           // TODO 需要配合 grpc.client 也需要写对于代码
	return srv
}

func WithAddress(address string) ServerOption {
	return func(s *Server) {
		s.address = address
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

func WithHealthServer(health *health.Server) ServerOption {
	return func(s *Server) {
		s.health = health
	}
}

func WithServerTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
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
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{
		Scheme: "grpc",
		Host:   addr,
	}
	return nil
}

// Start 启动 grpc 的服务
func (s *Server) Start() error {
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume() // 这里可以不写  只不过 我这里显示声明一下(健康检查服务知道该服务已经恢复正常工作)
	return s.Server.Serve(s.lis)
}

func (s *Server) Stop() error {
	s.health.Shutdown() // 设置 服务的状态为 not_serving, 阻止 接收新的请求 过来
	s.GracefulStop()    // grpc 优雅退出
	log.Infof("[gRPC] server stopping")
	return nil
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Endpoint() *url.URL {
	return s.endpoint
}
