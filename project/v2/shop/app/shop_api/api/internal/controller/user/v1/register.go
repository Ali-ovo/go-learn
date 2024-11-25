package controllerUser

import (
	"shop/gmicro/pkg/common/core"
	"shop/gmicro/pkg/log"
	translatorGin "shop/pkg/translator/gin"

	"github.com/gin-gonic/gin"
)

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻, 自定义validator
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}

func (us *userController) Register(ctx *gin.Context) {
	log.Info("Register is called")

	// 表单验证
	regForm := RegisterForm{}
	if err := ctx.ShouldBind(&regForm); err != nil {
		translatorGin.HandleValidatorError(ctx, err, us.trans)
		return
	}

	userDTO, err := us.srv.User().Register(ctx, regForm.Mobile, regForm.Password, regForm.Code)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, gin.H{
		"id":         userDTO.ID,
		"nick_name":  userDTO.NickName,
		"token":      userDTO.Token,
		"expired_at": userDTO.ExpiredAt,
	})
}
