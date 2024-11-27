package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type ServerOptions struct {
	EnableProfiling   bool     `json:"profiling"             mapstructure:"profiling"`           // 是否开启 pprof
	EnableLimit       bool     `json:"enable-limit"          mapstructure:"enable-limit"`        // 是否开启 限流设置
	EnableMetrics     bool     `json:"enable-metrics"        mapstructure:"enable-metrics"`      // 是否开启 metrics
	EnableHealthCheck bool     `json:"enable-health-check"   mapstructure:"enable-health-check"` // 是否开启 health check
	EnableTelemetry   bool     `json:"enable-telemetry"      mapstructure:"enable-telemetry"`    // 是否开启 telemetry 链路追踪
	Host              string   `json:"host,omitempty"        mapstructure:"host"`                // host
	Port              int      `json:"port,omitempty"        mapstructure:"port"`                // port
	HttpPort          int      `json:"http-port,omitempty"   mapstructure:"http-port"`           // http port
	Name              string   `json:"name,omitempty"        mapstructure:"name"`                // 名称
	HttpMode          string   `json:"http-mode,omitempty"   mapstructure:"http-mode"`           // gin 模式设置
	Middlewares       []string `json:"middlewares,omitempty" mapstructure:"middlewares"`         // Http中间件
}

// NewServerOptions creates a new ServerOptions
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		EnableProfiling: true,
		EnableMetrics:   true,
		EnableTelemetry: true,
		Host:            "127.0.0.1",
		Port:            8081,
		HttpPort:        8082,
		Name:            "Server",
		HttpMode:        "debug",
	}
}

// Validate verifies flags passed to ServerOptions
func (so *ServerOptions) Validate() []error {
	var errs []error
	if so.HttpPort == 0 || so.HttpPort > 65535 {
		errs = append(errs, fmt.Errorf("not a valid http port: %d", so.HttpPort))
	}

	if so.HttpPort != 0 {
		if !(so.HttpMode == "debug" || so.HttpMode == "test" || so.HttpMode == "release") {
			errs = append(errs, fmt.Errorf("mode must be one of debug/release/test"))
		}
	}
	return errs
}

// AddFlags adds flags related to server storage for a specific APIServer to specified FlagSet.
func (so *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&so.EnableProfiling, "server.profiling", so.EnableProfiling, "enable-profiling, if true, will add <host>/debug/pprof/ default is true")
	fs.BoolVar(&so.EnableLimit, "server.enable-limit", so.EnableLimit, "enable-limit")
	fs.BoolVar(&so.EnableMetrics, "server.enable-metrics", so.EnableMetrics, "enable-metrics, if true, will add /metrics, default is true")
	fs.BoolVar(&so.EnableHealthCheck, "server.enable-health-check", so.EnableHealthCheck, "enable-health-check, if true, will add health check route, default is true")
	fs.BoolVar(&so.EnableTelemetry, "server.enable-telemetry", so.EnableTelemetry, "enable-telemetry, if true, will add /telemetry, default is true")
	fs.StringVar(&so.Host, "server.host", so.Host, "server host default is 127.0.0.1")
	fs.IntVar(&so.Port, "server.port", so.Port, "server port default is 8081")
	fs.IntVar(&so.HttpPort, "server.http-port", so.HttpPort, "server http port default is 8082")
	fs.StringVar(&so.Name, "server.name", so.Name, "server name default is grpcServer")
	fs.StringVar(&so.HttpMode, "server.http-mode", so.HttpMode, "server http mode must be one of debug/release/test")
	fs.StringArrayVar(&so.Middlewares, "server.http-middlewares", so.Middlewares, "server middleware to be in recovery/cors/context")
}
