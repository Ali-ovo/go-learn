package main

import (
	"fmt"
	"go-learn/shop/shop_api/goods_web/global"
	"go-learn/shop/shop_api/goods_web/initialize"
	"go-learn/shop/shop_api/goods_web/utils"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	initialize.InitLogger()

	// 初始化配置
	initialize.InitConfig()

	// 初始化路由
	Router := initialize.Routers()

	// 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 初始化 srv 链接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	// 如果是本地开发环境端口号固定
	debug := viper.GetBool("DEBUG")
	if debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	zap.S().Infof("启动服务, 端口: %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
