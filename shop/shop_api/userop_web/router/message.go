package router

import (
	"go-learn/shop/shop_api/userop_web/api/message"
	"go-learn/shop/shop_api/userop_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 轮播图列表页
		MessageRouter.POST("", message.New) //新建轮播图
	}
}
