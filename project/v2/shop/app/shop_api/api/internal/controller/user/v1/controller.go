package controllerUser

import (
	"shop/app/shop_api/api/internal/service"

	ut "github.com/go-playground/universal-translator"
)

type userController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewUserController(trans ut.Translator, srv service.ServiceFactory) *userController {
	return &userController{trans: trans, srv: srv}
}
