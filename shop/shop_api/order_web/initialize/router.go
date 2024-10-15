package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-learn/shop/shop_api/order_web/middlewares"
	"go-learn/shop/shop_api/order_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// cors
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
