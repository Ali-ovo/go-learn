package router

import (
	"go-learn/shop/shop_api/user_web/api"
	"go-learn/shop/shop_api/user_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的 url")

	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}

}
