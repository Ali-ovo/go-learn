package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type ServerOptions struct {
	EnableProfiling   bool   `json:"profiling" mapstructure:"profiling"`
	EnableLimit       bool   `json:"limit" mapstructure:"limit"`                             // 是否开启pprof
	EnableMetrics     bool   `json:"enable-metrics" mapstructure:"enable-metrics"`           // 是否开启metrics
	EnableHealthCheck bool   `json:"enable-health-check" mapstructure:"enable-health-check"` // 是否开启health check
	Host              string `json:"host,omitempty" mapstructure:"host"`                     // host
	Port              int    `json:"port,omitempty" mapstructure:"port"`                     // port
	HttpPort          int    `json:"http-port,omitempty" mapstructure:"http-port"`           // http port
	Name              string `json:"name,omitempty" mapstructure:"name"`                     // 名称
	// TODO 需要处理
	HttpMiddlewares []string `json:"http-middlewares,omitempty" mapstructure:"http-middlewares"` // Http中间件
	GrpcMiddlewares []string `json:"grpc-middlewares,omitempty" mapstructure:"grpc-middlewares"` // grpc中间件
}

// NewServerOptions creates a new ServerOptions
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		EnableProfiling:   true,
		EnableMetrics:     true,
		EnableHealthCheck: true,
		Host:              "192.168.189.128",
		Port:              8081,
		HttpPort:          8082,
		Name:              "grpcServer",
	}
}

// Validate verifies flags passed to ServerOptions
func (so *ServerOptions) Validate() []error {
	var errs []error
	if so.HttpPort == 0 || so.HttpPort > 65535 {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", so.HttpPort))
	}
	return errs
}

// AddFlags adds flags related to server storage for a specific APIServer to specified FlagSet.
func (so *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&so.EnableProfiling, "server.profiling", so.EnableProfiling, "enable-profiling, if true, will add <host>/debug/pprof/ default is true")
	fs.BoolVar(&so.EnableMetrics, "server.enable-metrics", so.EnableMetrics, "enable-metrics, if true, will add /metrics, default is true")
	fs.BoolVar(&so.EnableHealthCheck, "server.enable-health-check", so.EnableHealthCheck, "enable-health-check, if true, will add health check route, default is true")
	fs.StringVar(&so.Host, "server.host", so.Host, "server host default is 127.0.0.1")
	fs.IntVar(&so.Port, "server.port", so.Port, "server port default is 8081")
	fs.IntVar(&so.HttpPort, "server.http-port", so.HttpPort, "server http port default is 8082")
	fs.StringVar(&so.Name, "server.name", so.Name, "server name default is grpcServer")
}
