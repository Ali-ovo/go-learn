package router

import (
	"go-learn/shop/shop_api/goods_web/api/category"
	"go-learn/shop/shop_api/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys").Use(middlewares.Trace())

	{
		CategoryRouter.GET("", category.List)                                                            // 商品列表
		CategoryRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.New)          // 新建商品
		CategoryRouter.GET("/:id", category.Detail)                                                      // 商品详情
		CategoryRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Delete) // 删除商品
		CategoryRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Update)    // 更新商品
	}

}
