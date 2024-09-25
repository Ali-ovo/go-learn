package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 绑定 JSON
type Login struct {
	User     string `json:"user" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required"`
}

type SignUpForm struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

var trans ut.Translator

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}

	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err

	}
	return rsp
}

func InitTrans(local string) (err error) {
	// 修改 gin 的 validator 引擎
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()

		// 第一个参数是备用语言，后面的参数是应该支持的语言
		uni := ut.New(
			enT,
			zhT,
		)

		trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", local)
		}

		switch local {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}

		return

	}

	return

}

func main() {
	if err := InitTrans("zh"); err != nil {
		fmt.Println("init translate error")

		return
	}

	r := gin.Default()

	r.POST("/login", func(ctx *gin.Context) {
		var login Login
		if err := ctx.ShouldBind(&login); err != nil {
			errs, ok := err.(validator.ValidationErrors)

			if !ok {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
			}

			ctx.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errs.Translate(trans))})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": login.User, "password": login.Password})
	})

	r.POST("./signup", func(ctx *gin.Context) {
		var SignUpForm SignUpForm
		if err := ctx.ShouldBind(&SignUpForm); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "sign in success"})
	})

	r.Run(":8081") // 监听并在 0.0.0.0:8080 上启动服务
}
