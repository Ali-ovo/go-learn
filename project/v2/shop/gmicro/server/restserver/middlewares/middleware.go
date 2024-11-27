package middlewares

import (
	sgin "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	"github.com/gin-gonic/gin"
)

var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),            // 将 在 错误 中恢复 并返回 500 http 状态码
		"cors":     Cors(),                    // 跨域访问
		"context":  Context(),                 // 自定义 中间件  在ctx 存储 ip 地址
		"limit":    sgin.SentinelMiddleware(), // 创建 熔断限流功能 拦截器
	}
}
