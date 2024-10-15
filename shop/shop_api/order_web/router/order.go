package router

import (
	"go-learn/shop/shop_api/order_web/api/order"
	"go-learn/shop/shop_api/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrdersRouter := Router.Group("orders").Use(middlewares.JWTAuth())

	{
		OrdersRouter.GET("", order.List)                            // 订单列表
		OrdersRouter.POST("", middlewares.IsAdminAuth(), order.New) // 新建订单
		OrdersRouter.GET("/:id", order.Detail)                      // 订单详情
	}
}
