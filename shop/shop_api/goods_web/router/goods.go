package router

import (
	"go-learn/shop/shop_api/goods_web/api/goods"
	"go-learn/shop/shop_api/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")

	{
		GoodsRouter.GET("", goods.List)                                                                 // 商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)               // 新建商品
		GoodsRouter.GET("/:id", goods.Detail)                                                           // 商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)      // 删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                                    // 商品库存
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) // 更新商品
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         // 更新商品
	}

}
