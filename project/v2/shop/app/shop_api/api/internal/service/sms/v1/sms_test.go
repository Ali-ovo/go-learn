package serviceSms

import (
	"context"
	"shop/pkg/options"
	"testing"
)

func TestSms(t *testing.T) {
	var smsOpts = options.SmsOptions{
		APIKey:      "LTAI5tGPjtFKjhq5Tstesm2d",
		APISecret:   "YlwEutdjZ8jCGBwBDaK7akgreSnKqe",
		APISignName: "生鲜电商",
		APICode:     "SMS_243216478",
	}

	s := NewSmsService(&smsOpts)
	err := s.SendSms(context.Background(), "13067353692", `{"code":"1234"}`)
	if err != nil {
		panic(err)
	}
}
