package admin

import (
	"shop/app/shop/api/config"
	"shop/app/shop/api/internal/controller/user/v1"
	"shop/app/shop/api/internal/data/v1/rpc"
	user2 "shop/app/shop/api/internal/service/user/v1"
	"shop/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	userGroup := v1.Group("/user")

	userData, err := rpc.GetDatafactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}

	userService := user2.NewUserService(userData, cfg.Jwt)
	uController := user.NewUserController(g.Translator(), userService)
	{
		userGroup.POST("pwd_login", uController.Login)
	}

}
