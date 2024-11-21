package user

import (
	"shop/app/shop/api/internal/service/user/v1"

	ut "github.com/go-playground/universal-translator"
)

type userController struct {
	trans ut.Translator
	srv   user.UserSrv
}

func NewUserController(trans ut.Translator, srv user.UserSrv) *userController {
	return &userController{trans: trans, srv: srv}
}
