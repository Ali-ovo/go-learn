package main

import (
	"fmt"

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

func main() {
	v := viper.New()

	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	serverConfig := ServerConfig{}

	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}

	fmt.Println(serverConfig)

}
