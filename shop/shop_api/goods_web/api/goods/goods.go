package goods

import (
	"context"
	"go-learn/shop/shop_api/goods_web/api"
	"go-learn/shop/shop_api/goods_web/forms"
	"go-learn/shop/shop_api/goods_web/global"
	"go-learn/shop/shop_api/goods_web/proto"
	"net/http"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 商品列表
func List(ctx *gin.Context) {

	request := &proto.GoodsFilterRequest{}
	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "1")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "10")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	// parent, _ := ctx.Get("parentSpan")
	// opentracing.ContextWithSpan(context.Background(), parent.(opentracing.Span))

	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {

		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过多",
		})
		return
	}

	rsp, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("获取商品列表失败", "msg", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()

	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}

	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {

	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// TODO 商品库存
	ctx.JSON(http.StatusOK, rsp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodInfoRequest{
		Id: int32(goodsId),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"goods_sn":    r.GoodsSn,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(ctx, &proto.DeleteGoodsInfo{
		Id: int32(goodsId),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	// TODO 商品库存

	return
}

func UpdateStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	goodsId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}

	if _, err = global.GoodsSrvClient.UpdateGoods(ctx, &proto.CreateGoodsInfo{
		Id:     int32(goodsId),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功"})

}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")

	goodsClient := global.GoodsSrvClient

	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	if _, err = goodsClient.UpdateGoods(ctx, &proto.CreateGoodsInfo{
		Id:              int32(goodsId),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
