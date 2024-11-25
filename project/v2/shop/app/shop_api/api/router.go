package api

import (
	"shop/app/shop_api/api/config"

	controllerSms "shop/app/shop_api/api/internal/controller/sms/v1"
	controllerUser "shop/app/shop_api/api/internal/controller/user/v1"
	"shop/app/shop_api/api/internal/data/v1/rpc"
	serviceSms "shop/app/shop_api/api/internal/service/sms/v1"
	serviceUser "shop/app/shop_api/api/internal/service/user/v1"
	"shop/app/shop_api/api/pkg/auth/BasicAuth"
	"shop/app/shop_api/api/pkg/auth/JWTAuth"
	"shop/gmicro/server/restserver"
	"shop/gmicro/server/restserver/middlewares/auth"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	userGroup := v1.Group("/user")
	baseGroup := v1.Group("/base")

	userData, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}

	userService := serviceUser.NewUserService(userData, cfg.Jwt)
	uController := controllerUser.NewUserController(g.Translator(), userService)
	// 做 jwt 校验用
	jwtAuth := JWTAuth.NewJWTAuth(cfg.Jwt)
	//jwtStragy := jwtAuth.(auth.JWTStrategy)
	//jwtStragy.LoginHandler
	// 做 basic 校验用
	basicAuth := BasicAuth.NewBasicAuth(userService)
	// jwt 和 basic 自动适配
	_ = auth.NewAutoStrategy(basicAuth, jwtAuth)
	{
		userGroup.POST("pwd_login", uController.Login)
		userGroup.POST("register", uController.Register)
		userGroup.GET("detail", jwtAuth.AuthFunc(), uController.GetUserDetail)
		userGroup.PATCH("update", jwtAuth.AuthFunc(), uController.UpdateUser)
	}

	smsService := serviceSms.NewSmsService(cfg.Sms)
	sController := controllerSms.NewSmsController(g.Translator(), smsService)
	{
		baseGroup.GET("captcha", controllerUser.GetCaptcha)
		baseGroup.POST("send_sms", sController.SendSms)
	}

}
