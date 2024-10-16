package pay

import (
	"fmt"
	"go-learn/shop/shop_api/order_web/global"
	"go-learn/shop/shop_api/order_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func Notify(ctx *gin.Context) {
	// 支付宝回调通知
	// 生成支付 url
	AliPayInfo := global.ServerConfig.AliPayInfo
	client, err := alipay.New(AliPayInfo.AppID, AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("创建支付宝客户端失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("获取支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	noti, err := client.GetTradeNotification(ctx.Request)
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
