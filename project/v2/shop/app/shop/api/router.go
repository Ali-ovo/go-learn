package api

import (
	"shop/app/shop/api/config"
	"shop/app/shop/api/internal/controller/sms/v1"
	"shop/app/shop/api/internal/controller/user/v1"
	"shop/app/shop/api/internal/data/v1/rpc"
	user2 "shop/app/shop/api/internal/service/user/v1"
	"shop/gmicro/server/restserver"

	sms2 "shop/app/shop/api/internal/service/sms/v1"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	userGroup := v1.Group("/user")
	baseGroup := v1.Group("/base")

	userData, err := rpc.GetDatafactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}

	userService := user2.NewUserService(userData, cfg.Jwt)
	uController := user.NewUserController(g.Translator(), userService)
	{
		userGroup.POST("pwd_login", uController.Login)
		userGroup.POST("register", uController.Register)
	}

	smsService := sms2.NewSmsService(cfg.Sms)
	sController := sms.NewSmsController(g.Translator(), smsService)
	{
		baseGroup.GET("captcha", user.GetCaptcha)
		baseGroup.POST("send_sms", sController.SendSms)
	}

}
