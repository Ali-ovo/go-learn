package restserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	mws "shop/gmicro/server/restserver/middlewares"
	"shop/gmicro/server/restserver/pprof"
	"shop/gmicro/server/restserver/validation"
	"net/http"
	"time"
)

type JwtInfo struct {
	// default to "JWT"
	Realm string
	// default to empty
	Key string
	// default to 7 days
	Timeout time.Duration
	// default to 7 days
	MaxRefresh time.Duration
}

// Server wrapper for gin.Engine
type Server struct {
	*gin.Engine
	// 端口号
	port int
	// 开发模式 默认值 debug
	mode string
	// 是否开启健康检查接口 默认开启, 如果开启会自动添加 /health
	healthz bool
	// 是否开启 pprof接口, 默认开启, 如果开启会自动添加 /debug/pprof 接口
	enableProfiling bool
	// 是否开启 链路追踪 默认开启
	enableTracing bool
	// 中间件
	//customMiddlewares []gin.HandlerFunc
	middlewares []string // 这里不选择注入的方式  选择做好的方法 选择用即可
	//jwt配置信息
	jwt *JwtInfo
	// 翻译器
	transName   string
	trans       ut.Translator
	server      *http.Server
	serviceName string
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		port:            8080,
		mode:            "debug",
		healthz:         true,
		enableProfiling: true,
		enableTracing:   true,
		jwt: &JwtInfo{
			Realm:      "JWT",
			Key:        "3KcjkZbUDdEeGeFX@h^!qKXh2WC@A6Qe",
			Timeout:    time.Hour * 24 * 7,
			MaxRefresh: time.Hour * 24 * 7,
		},
		Engine:      gin.Default(),
		transName:   "zh",
		serviceName: "gmicro",
	}

	for _, o := range opts {
		o(srv)
	}

	if srv.enableTracing {
		srv.Use(mws.TracingHandler(srv.serviceName))
	}

	for _, m := range srv.middlewares {
		mw, ok := mws.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
			//panic(errors.Errorf("can not find middleware: %s", m))
		}

		log.Infof("intall middleware: %s", m)
		srv.Use(mw)
	}

	return srv
}

// Start rest server
func (s *Server) Start(ctx context.Context) error {
	// 设置开发模式, 打印路由信息
	if s.mode != gin.DebugMode && s.mode != gin.ReleaseMode && s.mode != gin.TestMode {
		return errors.New("mode must be one of debug/release/test")
	}
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	err := s.initTrans(s.transName)
	if err != nil {
		log.Errorf("init translator error: %s", err.Error())
		return err
	}

	// 注册 mobile 验证器
	validation.RegisterMobile(s.trans)

	// 根据配置初始化 pprof 路由
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	log.Infof("start rest server on port: %d", s.port)
	address := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}
	// 在 Gin 框架中，SetTrustedProxies 方法是 gin.Engine 结构体的一个方法，用于设置信任的代理服务器，
	// 以便在处理 HTTP 请求时跳过这些代理服务器的 IP 地址。将参数设置为 nil 表示没有信任的代理服务器需要跳过，这可能会影响请求的处理方式
	_ = s.SetTrustedProxies(nil)
	//err = s.Run(fmt.Sprintf(":%d", s.port))
	// http.ErrServerClosed 表示服务已经优雅退出了
	if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("stop rest server")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("shutdown rest server error: %s", err.Error())
		return err
	}
	log.Info("rest server stopped")
	return nil
}
