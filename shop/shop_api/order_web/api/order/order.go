package order

import (
	"go-learn/shop/shop_api/order_web/api"
	"go-learn/shop/shop_api/order_web/global"
	"go-learn/shop/shop_api/order_web/models"
	"go-learn/shop/shop_api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	// 获取订单列表
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)

	request := proto.OrderFilterRequest{}

	// admin
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	request.Pages = int32(pnInt)
	request.PagePerNums = int32(pSizeInt)
	rsp, err := global.OrderSrvClient.OrderList(ctx, &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		orderList = append(orderList, map[string]interface{}{
			"id":       item.Id,
			"status":   item.Status,
			"pay_type": item.PayType,
			"user":     item.UserId,
			"post":     item.Post,
			"total":    item.Total,
			"address":  item.Address,
			"name":     item.Name,
			"mobile":   item.Mobile,
			"order_sn": item.OrderSn,
			"add_time": item.AddTime,
		})
	}

	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)

}

func New(ctx *gin.Context) {

}

func Detail(ctx *gin.Context) {

}
