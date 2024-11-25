package controllerUser

import (
	"net/http"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/common/core"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	translatorGin "shop/pkg/translator/gin"

	"github.com/gin-gonic/gin"
)

type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻， 自定义validate
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=7"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

func (us *userController) Login(ctx *gin.Context) {
	log.Info("Login is called")

	// 表单验证
	passwordLoginForm := PassWordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		translatorGin.HandleValidatorError(ctx, err, us.trans)
		return
	}

	// TODO 先关闭着
	//// 图片验证码验证
	//if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
	//	core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, "验证码错误"), nil)
	//	return
	//}
	userDTO, err := us.srv.User().MobileLogin(ctx, passwordLoginForm.Mobile, passwordLoginForm.PassWord)
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPasswordIncorrect, "登录失败"), nil)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":         userDTO.ID,
		"nick_name":  userDTO.NickName,
		"token":      userDTO.Token,
		"expired_at": userDTO.ExpiredAt,
	})
}
