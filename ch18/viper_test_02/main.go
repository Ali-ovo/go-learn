package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type MySqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	MysqlInfo MySqlConfig `mapstructure:"mysql"`
	Name      string      `mapstructure:"name"`
}

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func main() {
	debug := GetEnvInfo("DEBUG")
	fmt.Println(debug)
	configFileName := "config_dev.yaml"
	if debug == "true" {
		configFileName = "config_pro.yaml"
	}
	v := viper.New()

	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	serverConfig := ServerConfig{}

	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		_ = v.ReadInConfig()
		if err := v.Unmarshal(&serverConfig); err != nil {
			panic(err)
		}

		fmt.Println(serverConfig)
	})

}
