package shopcart

import (
	"go-learn/shop/shop_api/order_web/api"
	"go-learn/shop/shop_api/order_web/forms"
	"go-learn/shop/shop_api/order_web/global"
	"go-learn/shop/shop_api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	// 获取购物车商品
	userId, _ := ctx.Get("userId")

	rsp, err := global.OrderSrvClient.CartItemList(ctx, &proto.UserInfo{
		Id: int32(userId.(uint)),
	})

	if err != nil {
		zap.S().Errorw("查询购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 获取商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(ctx, &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("查询商品信息失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		for _, goods := range goodsRsp.Data {
			if item.GoodsId == goods.Id {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["goods_name"] = goods.Name
				tmpMap["goods_price"] = goods.ShopPrice
				tmpMap["goods_image"] = goods.GoodsFrontImage
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}

	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	// 添加商品到购物
	itemForm := forms.ShopCartItemForm{}

	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 添加商品之前检查商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(ctx, &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})

	if err != nil {
		zap.S().Errorw("查询商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	invRsp, err := global.InventorySrvClient.InvDetail(ctx, &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查询库存失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")

	rsp, err := global.OrderSrvClient.CreateCartItem(ctx, &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})

	if err != nil {
		zap.S().Errorw("添加购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "请求参数错误",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(ctx, &proto.CartItemRequest{
		GoodsId: int32(i),
		UserId:  int32(userId.(uint)),
	})

	if err != nil {
		zap.S().Errorw("删除购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	// 更新购物车商品
	itemForm := forms.ShopCartItemUpdateForm{}
	if er := ctx.ShouldBindJSON(&itemForm); er != nil {
		api.HandleValidatorError(ctx, er)
		return
	}

	userId, _ := ctx.Get("userId")
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "请求参数错误",
		})
		return
	}

	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(ctx, &request)

	if err != nil {
		zap.S().Errorw("更新购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
