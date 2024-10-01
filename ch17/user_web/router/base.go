package router

import (
	"go-learn/ch17/user_web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")

	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
		BaseRouter.POST("/send_sms", api.SendSms)
	}
}
