package restserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"shop/gmicro/pkg/common/endpoint"
	"shop/gmicro/pkg/host"
	"shop/gmicro/pkg/log"
	"shop/gmicro/server/restserver/middlewares"
	"shop/gmicro/server/restserver/pprof"
	"shop/gmicro/server/restserver/validation"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

//type JwtInfo struct {
//	// default to "JWT"
//	Realm string
//	// default to empty
//	Key string
//	// default to 7 days
//	Timeout time.Duration
//	// default to 7 days
//	MaxRefresh time.Duration
//}

// Server wrapper for gin.Engine
type Server struct {
	*gin.Engine
	address  string
	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL
	err      error
	// 开发模式 默认值 debug
	mode string
	// 是否开启 pprof接口, 默认开启, 如果开启会自动添加 /debug/pprof 接口
	enableProfiling bool
	// 是否开启 链路追踪 默认开启
	enableTracing bool
	// 是否开启 metrics 接口, 默认开启, 如果开启会自动添加 /metrics 接口
	enableMetrics bool
	// 中间件
	//customMiddlewares []gin.HandlerFunc
	middlewares []string // 这里不选择注入的方式  选择做好的方法 选择用即可
	////jwt配置信息
	//jwt *JwtInfo
	// 翻译器
	transName   string
	trans       ut.Translator
	server      *http.Server
	serviceName string
}

func (s *Server) Translator() ut.Translator {
	return s.trans
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		address:         ":0",
		mode:            "debug",
		enableProfiling: true,
		enableTracing:   true,
		Engine:          gin.Default(),
		transName:       "zh",
		serviceName:     "gmicro",
	}

	for _, o := range opts {
		o(srv)
	}

	// 开启链路追踪
	if srv.enableTracing {
		srv.Use(middlewares.TracingHandler(srv.serviceName))
	}

	for _, m := range srv.middlewares {
		mw, ok := middlewares.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
			//panic(errors.Errorf("can not find middleware: %s", m))
		}

		log.Infof("intall middleware: %s", m)
		srv.Use(mw)
	}

	// 设置开发模式, 打印路由信息
	if srv.mode != gin.DebugMode && srv.mode != gin.ReleaseMode && srv.mode != gin.TestMode {
		panic("mode must be one of debug/release/test")
	}
	gin.SetMode(srv.mode)
	// 修改 gin 调试日志输出格式
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		handlerName = strings.SplitN(handlerName, ".", 2)[0]
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	err := srv.initTrans(srv.transName)
	if err != nil {
		panic(fmt.Sprintf("init translator error: %s", err.Error()))
	}

	// 注册 mobile 验证器
	validation.RegisterMobile(srv.trans)

	// 根据配置初始化 pprof 路由
	if srv.enableProfiling {
		pprof.Register(srv.Engine) // TODO 不知道有什么用
	}

	// 开启 普罗米修斯监控
	if srv.enableMetrics { // TODO 需要重新巩固一下
		m := ginmetrics.GetMonitor() // gin-metrics 可以自动 收集相关信息 并且开启 /metrics 的 接口

		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// 慢请求的时间, default 5s
		m.SetSlowTime(5)
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		// set middleware for gin
		m.Use(srv.Engine)
	}

	return srv
}

// Start rest server
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[gin] start rest server on %s", s.address)
	s.server = &http.Server{
		Addr:    s.address,
		Handler: s.Engine,
	}
	// 在 Gin 框架中，SetTrustedProxies 方法是 gin.Engine 结构体的一个方法，用于设置信任的代理服务器，
	// 以便在处理 HTTP 请求时跳过这些代理服务器的 IP 地址。将参数设置为 nil 表示没有信任的代理服务器需要跳过，这可能会影响请求的处理方式
	_ = s.SetTrustedProxies(nil)
	//err = s.Run(fmt.Sprintf(":%d", s.port))
	// http.ErrServerClosed 表示服务已经优雅退出了
	var err error
	if s.tlsConf != nil {
		// 需要你填写具体的证书
		err = s.server.ServeTLS(s.lis, "", "")
	} else {
		err = s.server.Serve(s.lis)
	}
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("shutdown rest server error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

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
			_ = s.lis.Close()
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.tlsConf != nil), addr)
	}
	return s.err
}
