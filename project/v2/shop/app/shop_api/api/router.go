package api

import (
	"shop/app/shop_api/api/config"
	"shop/app/shop_api/api/controller/sms/v1"
	"shop/app/shop_api/api/controller/user/v1"
	"shop/app/shop_api/api/data/v1/rpc"
	"shop/app/shop_api/api/pkg/auth/BasicAuth"
	"shop/app/shop_api/api/pkg/auth/JWTAuth"
	sms2 "shop/app/shop_api/api/service/sms/v1"
	user2 "shop/app/shop_api/api/service/user/v1"
	"shop/gmicro/server/restserver"
	"shop/gmicro/server/restserver/middlewares/auth"
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

	smsService := sms2.NewSmsService(cfg.Sms)
	sController := sms.NewSmsController(g.Translator(), smsService)
	{
		baseGroup.GET("captcha", user.GetCaptcha)
		baseGroup.POST("send_sms", sController.SendSms)
	}

}
