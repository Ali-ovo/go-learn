package initialize

import (
	"github.com/gin-gonic/gin"

	"go-learn/ch17/user_web/middlewares"
	"go-learn/ch17/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// cors
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return Router
}
