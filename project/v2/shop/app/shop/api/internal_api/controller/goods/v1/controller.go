package controllerGoods

import (
	"shop/app/shop/api/internal_api/service"

	ut "github.com/go-playground/universal-translator"
)

type GoodsController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewGoodsController(trans ut.Translator, srv service.ServiceFactory) *GoodsController {
	return &GoodsController{
		trans: trans,
		srv:   srv,
	}
}
