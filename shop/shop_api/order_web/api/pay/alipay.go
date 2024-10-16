package pay

import (
	"fmt"
	"go-learn/shop/shop_api/order_web/global"
	"go-learn/shop/shop_api/order_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Notify(ctx *gin.Context) {
	// 支付宝回调通知

	noti, err := global.AliPayClient.GetTradeNotification(ctx.Request)
	if err != nil {
		fmt.Println("交易状态为：", noti.TradeStatus)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	_, err = global.OrderSrvClient.UpdateOrderStatus(ctx, &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.String(http.StatusOK, "success")
}
