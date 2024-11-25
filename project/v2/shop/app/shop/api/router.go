package api

import (
	"shop/app/shop/api/config"
	controllerGoods "shop/app/shop/api/internal_api/controller/goods/v1"
	controllerSms "shop/app/shop/api/internal_api/controller/sms/v1"
	controllerUser "shop/app/shop/api/internal_api/controller/user/v1"
	"shop/app/shop/api/internal_api/data/v1/rpc"
	"shop/app/shop/api/internal_api/service"
	"shop/app/shop/api/pkg/auth/BasicAuth"

	"shop/app/shop_api/api/pkg/auth/JWTAuth"
	"shop/gmicro/server/restserver"
	"shop/gmicro/server/restserver/middlewares/auth"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	baseGroup := v1.Group("/base")
	userGroup := v1.Group("/user")
	goodsGroup := v1.Group("/goods")

	data, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}
	serviceFactory := service.NewService(data, cfg.Sms, cfg.Jwt)

	uController := controllerUser.NewUserController(g.Translator(), serviceFactory)
	jwtAuth := JWTAuth.NewJWTAuth(cfg.Jwt) // 做 jwt 校验用
	//jwtStragy := jwtAuth.(auth.JWTStrategy)
	//jwtStragy.LoginHandler
	basicAuth := BasicAuth.NewBasicAuth(serviceFactory.User()) // 做 basic 校验用
	_ = auth.NewAutoStrategy(basicAuth, jwtAuth)               // jwt 和 basic 自动适配
	{
		userGroup.POST("pwd_login", uController.Login)
		userGroup.POST("register", uController.Register)
		userGroup.GET("detail", jwtAuth.AuthFunc(), uController.GetUserDetail)
		userGroup.PATCH("update", jwtAuth.AuthFunc(), uController.UpdateUser)
	}

	sController := controllerSms.NewSmsController(g.Translator(), serviceFactory)
	{
		baseGroup.GET("captcha", controllerUser.GetCaptcha)
		baseGroup.POST("send_sms", sController.SendSms)
	}

	gController := controllerGoods.NewGoodsController(g.Translator(), serviceFactory)
	{
		goodsGroup.GET("", gController.List)
	}
}
