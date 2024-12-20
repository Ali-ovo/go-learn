package controllerGoods

import (
	"shop/app/shop_api/api/internal/service"

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
