package config

import (
	cliflag "shop/gmicro/pkg/common/cli/flag"
	"shop/gmicro/pkg/log"
	"shop/pkg/options"
)

type Config struct {
	Log       *log.Options              `json:"log" mapstructure:"log"`             // log日志包 相关配置
	Server    *options.ServerOptions    `json:"server" mapstructure:"server"`       // server 相关配置
	Registry  *options.RegistryOptions  `json:"registry" mapstructure:"registry"`   // 服务注册发现注销 相关配置
	Jwt       *options.JwtOptions       `json:"jwt" mapstructure:"jwt"`             // jwt 配置参数
	Sms       *options.SmsOptions       `json:"sms" mapstructure:"sms"`             // 短信 相关配置
	Redis     *options.RedisOptions     `json:"redis" mapstructure:"redis"`         // redis 相关配置
	Telemetry *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"` // 链路追踪 相关配置
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
	errors = append(errors, c.Jwt.Validate()...)
	errors = append(errors, c.Sms.Validate()...)
	errors = append(errors, c.Redis.Validate()...)
	errors = append(errors, c.Telemetry.Validate()...)
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
	c.Jwt.AddFlags(fss.FlagSet("jwt"))
	c.Sms.AddFlags(fss.FlagSet("sms"))
	c.Redis.AddFlags(fss.FlagSet("redis"))
	c.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	return fss
}

func NewConfig() *Config {
	return &Config{
		Log:       log.NewOptions(),              // 初始化 log 配置
		Server:    options.NewServerOptions(),    // 初始化 Server 配置
		Registry:  options.NewRegistryOptions(),  // 初始化 consul 配置
		Jwt:       options.NewJwtOptions(),       // 初始化 jwt 配置
		Sms:       options.NewSmsOptions(),       // 初始化 sms 配置
		Redis:     options.NewRedisOptions(),     // 初始化 redis 配置
		Telemetry: options.NewTelemetryOptions(), // 初始化 telemetry 配置
	}
}
