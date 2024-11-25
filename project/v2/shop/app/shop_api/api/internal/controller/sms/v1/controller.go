package controllerSms

import (
	"shop/app/shop_api/api/internal/service"

	ut "github.com/go-playground/universal-translator"
)

type smsController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewSmsController(trans ut.Translator, sms service.ServiceFactory) *smsController {
	return &smsController{
		trans: trans,
		srv:   sms,
	}
}
