package restserver

import "crypto/tls"

type ServerOption func(*Server)

func WithEnableProfiling(enable bool) ServerOption {
	return func(s *Server) {
		s.enableProfiling = enable
	}
}

func WithMode(mode string) ServerOption {
	return func(s *Server) {
		s.mode = mode
	}
}

func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func WithTransNames(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}

// WithClientEnableTracing 设置是否开启链路追踪
func WithClientEnableTracing(enableTracing bool) ServerOption {
	return func(s *Server) {
		s.enableTracing = enableTracing
	}
}

func WithTransName(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}

func WithEnableMetrics(enable bool) ServerOption {
	return func(s *Server) {
		s.enableMetrics = enable
	}
}
