package sms

import (
	"shop/app/shop_api/api/service/sms/v1"
	"shop/gmicro/pkg/common/core"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/storage"
	"shop/pkg/code"
	gin2 "shop/pkg/translator/gin"
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type SendSmsForm struct {
	// 注册发送短信验证码和动态验证码登录发送验证码
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}

type smsController struct {
	trans ut.Translator
	srv   sms.SmsSrv
}

func NewSmsController(trans ut.Translator, sms sms.SmsSrv) *smsController {
	return &smsController{
		trans: trans,
		srv:   sms,
	}
}

func (ss *smsController) SendSms(ctx *gin.Context) {
	sendSmsForm := SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		gin2.HandleValidatorError(ctx, err, ss.trans)
		return
	}

	smsCode := sms.GenerateSmsCode(6)
	err := ss.srv.SendSms(ctx, sendSmsForm.Mobile, "{\"code\":"+smsCode+"}")
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	//将验证码保存起来 - redis
	rstore := storage.RedisCluster{}
	err = rstore.SetKey(ctx, sendSmsForm.Mobile, smsCode, 5*time.Minute)
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
