package initialize

import (
	"go-learn/shop/shop_api/order_web/global"

	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func InitAliPayClient() {
	// 支付宝回调通知
	// 生成支付 url
	AliPayInfo := global.ServerConfig.AliPayInfo
	var err error
	global.AliPayClient, err = alipay.New(AliPayInfo.AppID, AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("创建支付宝客户端失败")
		panic(err)
	}

	err = global.AliPayClient.LoadAliPayPublicKey(AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("获取支付宝公钥失败")
		panic(err)
	}

}
