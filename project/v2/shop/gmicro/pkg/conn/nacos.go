package conn

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

type NacosOptions struct {
	Host                 string
	Port                 uint64
	User                 string
	Password             string
	TimeOut              uint64
	NotLoadCacheAtStart  bool
	UpdateCacheWhenEmpty bool
	LogDir               string
	LogLevel             string
	CacheDir             string
	DataID               string
	Group                string
}

func NewNacosClient(opts *NacosOptions) (config_client.IConfigClient, error) {
	// nacos server地址
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			opts.Host,
			opts.Port,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}

	var clientOpts []constant.ClientOption
	if opts.User != "" {
		clientOpts = append(clientOpts, constant.WithUsername(opts.User), constant.WithPassword(opts.Password))
	}
	if opts.TimeOut != 0 {
		clientOpts = append(clientOpts, constant.WithTimeoutMs(opts.TimeOut))
	}
	if opts.NotLoadCacheAtStart {
		clientOpts = append(clientOpts, constant.WithNotLoadCacheAtStart(opts.NotLoadCacheAtStart))
	}
	if opts.UpdateCacheWhenEmpty {
		clientOpts = append(clientOpts, constant.WithUpdateCacheWhenEmpty(opts.UpdateCacheWhenEmpty))
	}
	if opts.LogDir != "" {
		clientOpts = append(clientOpts, constant.WithLogDir(opts.LogDir))
	}
	if opts.LogLevel != "" {
		clientOpts = append(clientOpts, constant.WithLogLevel(opts.LogLevel))
	}
	if opts.CacheDir != "" {
		clientOpts = append(clientOpts, constant.WithCacheDir(opts.CacheDir))
	}
	clientConfig := *constant.NewClientConfig(clientOpts...)

	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
}
