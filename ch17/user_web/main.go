package main

import (
	"fmt"
	"go-learn/ch17/user_web/initialize"

	"go.uber.org/zap"
)

func main() {
	port := 8090

	// 初始化日志
	initialize.InitLogger()

	// 初始化路由
	Router := initialize.Routers()

	zap.S().Infof("启动服务, 端口: %d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
