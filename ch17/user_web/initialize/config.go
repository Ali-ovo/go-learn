package initialize

import (
	"go-learn/ch17/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("DEBUG")
	configFileName := "config_dev.yaml"
	if debug {
		configFileName = "config_pro.yaml"
	}
	v := viper.New()

	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {

		zap.S().Infof("配置文件修改:%v", in.Name)
		_ = v.ReadInConfig()
		if err := v.Unmarshal(global.ServerConfig); err != nil {
			panic(err)
		}

		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})
}
