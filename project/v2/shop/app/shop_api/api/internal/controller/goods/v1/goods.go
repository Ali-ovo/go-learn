package controllerGoods

import (
	"net/http"
	dtoGoods "shop/app/shop_api/api/internal/domain/dto/goods"
	"shop/gmicro/pkg/common/core"
	translatorGin "shop/pkg/translator/gin"

	"github.com/gin-gonic/gin"
)

func (gc *GoodsController) List(ctx *gin.Context) {
	var r dtoGoods.GoodsFilter

	if err := ctx.ShouldBindQuery(&r); err != nil {
		translatorGin.HandleValidatorError(ctx, err, gc.trans)
		return
	}

	goodsDTOList, err := gc.srv.Goods().GoodsList(ctx, &r)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	ctx.JSON(http.StatusOK, goodsDTOList)
	return
}
