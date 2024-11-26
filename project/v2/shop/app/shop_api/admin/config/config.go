package config

import (
	cliflag "shop/gmicro/pkg/common/cli/flag"
	"shop/gmicro/pkg/log"
	options2 "shop/pkg/options"
)

type Config struct {
	Log      *log.Options              `json:"log" mapstructure:"log"`           // log 配置参数
	Server   *options2.ServerOptions   `json:"server" mapstructure:"server"`     // server
	Registry *options2.RegistryOptions `json:"registry" mapstructure:"registry"` // 服务注册 发现 注销 配置参数
}

func (c *Config) Validate() []error {
	var errors []error
	errors = append(errors, c.Log.Validate()...)
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Registry.Validate()...)
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
	return fss
}

func NewConfig() *Config {
	return &Config{
		Log:      log.NewOptions(),              // 初始化 log 配置
		Server:   options2.NewServerOptions(),   // 初始化 Server 配置
		Registry: options2.NewRegistryOptions(), // 初始化 consul 配置
	}
}
