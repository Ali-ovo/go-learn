package restserver

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

func (s *Server) initTrans(locale string) (err error) {
	//修改gin框架中的validator引擎属性, 实现定制
	// 详情参考 https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//v.RegisterValidation
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		s.trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, s.trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, s.trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, s.trans)
		}
		return
	}
	return
}
