package initialize

import (
	"encoding/json"
	"fmt"
	"go-learn/shop/shop_srvs/goods_srv/global"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("Ali_DEBUG")
	configFileName := "config_dev.yaml"
	if debug {
		configFileName = "config_pro.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConfig.Namespace), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("tmp/nacos/log"),
		constant.WithCacheDir("tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithUsername("nacos"),
		constant.WithPassword("nacos"),
	)

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			global.NacosConfig.Host,
			global.NacosConfig.Port,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)

	if err != nil {
		zap.S().Fatalf("unmarshal nacos content failed:%v", err)
	}
	fmt.Println(&global.ServerConfig)
}
