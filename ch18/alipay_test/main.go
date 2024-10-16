package main

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

func main() {
	appID := ""

	privateKey := ""

	aliPublicKey := ""

	var client, err = alipay.New(appID, privateKey, false)

	if err != nil {
		panic(err)
	}

	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://4nj7qikr.dongtaiyuming.net/o/v1/pay/alipay/notify"
	p.ReturnURL = "https://4nj7qikr.dongtaiyuming.net/o/v1/pay/alipay/return"
	p.Subject = "订单支付"
	p.OutTradeNo = "ali_go"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(url.String())

}
