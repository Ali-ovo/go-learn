package config

import (
	cliflag "shop/gmicro/pkg/common/cli/flag"
	"shop/gmicro/pkg/log"
	"shop/pkg/options"
)

type Config struct {
	Log       *log.Options              `json:"log"       mapstructure:"log"`       // log日志包 相关配置
	Server    *options.ServerOptions    `json:"server"    mapstructure:"server"`    // server 相关配置
	Registry  *options.RegistryOptions  `json:"registry"  mapstructure:"registry"`  // 服务注册发现注销 相关配置
	Telemetry *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"` // 链路追踪 相关配置
	Mysql     *options.MySQLOptions     `json:"mysql"     mapstructure:"mysql"`     // mysql 相关配置
	EsOptions *options.EsOptions        `json:"es"        mapstructure:"es"`        // es 相关配置
	Rocketmq  *options.RocketmqOptions  `json:"rocketmq"  mapstructure:"rocketmq"`  // rocketmq 相关配置
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
	errors = append(errors, c.Telemetry.Validate()...)
	errors = append(errors, c.Mysql.Validate()...)
	errors = append(errors, c.EsOptions.Validate()...)
	errors = append(errors, c.Rocketmq.Validate()...)
	return errors
}

// Flags
// fss cliflag.NamedFlagSets 声明一个 空结构体
// fss.FlagSet("log") 填充具体数据
func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	// fss.FlagSet("log") 生成 pflag 对象 还未解析数据
	// c.Log.AddFlags 设置需要解析的数据
	c.Log.AddFlags(fss.FlagSet("log"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	c.Mysql.AddFlags(fss.FlagSet("mysql"))
	c.EsOptions.AddFlags(fss.FlagSet("es"))
	c.Rocketmq.AddFlags(fss.FlagSet("rocketmq"))
	return fss
}

func NewConfig() *Config {
	return &Config{
		Log:       log.NewOptions(),              // 初始化 log 配置
		Server:    options.NewServerOptions(),    // 初始化 Server 配置
		Registry:  options.NewRegistryOptions(),  // 初始化 consul 配置
		Telemetry: options.NewTelemetryOptions(), // 初始化 telemetry 配置
		Mysql:     options.NewMySQLOptions(),     // 初始化 mysql 配置
		EsOptions: options.NewEsOptions(),        // 初始化 es 配置
		Rocketmq:  options.NewRocketmqOptions(),  // 初始化 rocketmq 配置
	}
}
