package order

import (
	"fmt"
	"go-learn/shop/shop_api/order_web/api"
	"go-learn/shop/shop_api/order_web/forms"
	"go-learn/shop/shop_api/order_web/global"
	"go-learn/shop/shop_api/order_web/models"
	"go-learn/shop/shop_api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
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
	orderForm := forms.CreateOrderForm{}

	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(ctx, &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})

	if err != nil {
		zap.S().Errorw("创建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成支付 url
	AliPayInfo := global.ServerConfig.AliPayInfo

	var p = alipay.TradePagePay{}
	p.NotifyURL = AliPayInfo.NotifyURL
	p.ReturnURL = AliPayInfo.ReturnURL
	p.Subject = rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := global.AliPayClient.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付 url 失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	fmt.Println(url.String())

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})

}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	userId, _ := ctx.Get("userId")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "参数错误",
		})
		return
	}

	request := proto.OrderRequest{
		Id: int32(orderId),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)

	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(ctx, &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"id":       rsp.OrderInfo.Id,
		"status":   rsp.OrderInfo.Status,
		"user":     rsp.OrderInfo.UserId,
		"post":     rsp.OrderInfo.Post,
		"total":    rsp.OrderInfo.Total,
		"address":  rsp.OrderInfo.Address,
		"name":     rsp.OrderInfo.Name,
		"mobile":   rsp.OrderInfo.Mobile,
		"pay_type": rsp.OrderInfo.PayType,
		"order_sn": rsp.OrderInfo.OrderSn,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		goodsList = append(goodsList, map[string]interface{}{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"price": item.GoodsPrice,
			"image": item.GoodsImage,
			"nums":  item.Nums,
		})
	}

	reMap["goods"] = goodsList

	// 生成支付 url
	AliPayInfo := global.ServerConfig.AliPayInfo

	var p = alipay.TradePagePay{}
	p.NotifyURL = AliPayInfo.NotifyURL
	p.ReturnURL = AliPayInfo.ReturnURL
	p.Subject = rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := global.AliPayClient.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付 url 失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	fmt.Println(url.String())
	reMap["alipay_url"] = url.String()

	ctx.JSON(http.StatusOK, reMap)
}
