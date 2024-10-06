package router

import (
	"go-learn/shop/shop_api/user_web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")

	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
		BaseRouter.POST("/send_sms", api.SendSms)
	}
}
