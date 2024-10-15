package router

import (
	shopCart "go-learn/shop/shop_api/order_web/api/shop_cart"
	"go-learn/shop/shop_api/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())

	{
		ShopCartRouter.GET("", shopCart.List)          // 购物车列表
		ShopCartRouter.DELETE("/:id", shopCart.Delete) // 删除条目
		ShopCartRouter.POST("", shopCart.New)          // 新增条目
		ShopCartRouter.PATCH("/:id", shopCart.Update)  // 更新条目
	}
}
