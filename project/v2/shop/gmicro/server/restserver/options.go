package restserver

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

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMiddlewares(middlewares []string) ServerOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithHealthz(healthz bool) ServerOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithJwt(jwt *JwtInfo) ServerOption {
	return func(s *Server) {
		s.jwt = jwt
	}
}

func WithTransNames(transName string) ServerOption {
	return func(s *Server) {
		s.transName = transName
	}
}

func WithEnableTracing(enableTracing bool) ServerOption {
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
