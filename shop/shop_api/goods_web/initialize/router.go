package initialize

import (
	"github.com/gin-gonic/gin"

	"go-learn/shop/shop_api/goods_web/middlewares"
	"go-learn/shop/shop_api/goods_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// cors
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)

	return Router
}
