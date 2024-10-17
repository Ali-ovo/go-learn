package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-learn/shop/shop_api/userop_web/middlewares"
	"go-learn/shop/shop_api/userop_web/router"
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

	ApiGroup := Router.Group("/up/v1")
	router.InitUserFavRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)

	return Router
}
