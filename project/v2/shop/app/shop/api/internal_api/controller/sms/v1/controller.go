package controllerSms

import (
	"shop/app/shop/api/internal_api/service"

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
