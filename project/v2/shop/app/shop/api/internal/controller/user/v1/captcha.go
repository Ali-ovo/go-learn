package user

import (
	"image/color"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// 有 10 分钟超时机制 也可自己定义
var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	ctx.Query("captcha")

	// https://captcha.mojotv.cn/
	var driver base64Captcha.Driver
	switch ctx.Query("captcha") {
	case "audio":
		//driver = base64Captcha.NewDriverAudio()
	case "string":
		driver = base64Captcha.NewDriverString(
			80,                                     // 图片的 高度
			240,                                    // 图片的 宽度
			0,                                      // 背景干扰的字母
			0,                                      // 选择 背景的干扰的线条
			6,                                      // 字符的长度
			"1234567890qwertyuioplkjhgfdsazxcvbnm", // 选择出现的元素
			&color.RGBA{ // 背景颜色
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			},
			nil,
			[]string{"3Dumb.ttf"}, // 选择的字体
		)
	case "math":
		driver = base64Captcha.NewDriverMath(
			80,  // 图片的 高度
			240, // 图片的 宽度
			0,   // 背景干扰的字母
			2,   // 选择 背景的干扰的线条
			&color.RGBA{ // 背景颜色
				R: 20,
				G: 20,
				B: 20,
				A: 254,
			},
			nil,
			/*
				"3Dumb.ttf"
				"ApothecaryFont.ttf"
				"Flim-Flam.ttf"
				"Comismsh.ttf"
				"RitaSmith.ttf"
				"DENNEthree-dee.ttf"
				"actionj.ttf"
				"chromohv.ttf"
				"DeborahFancyDress.ttf"
				"wqy-microhei.ttc"
			*/
			[]string{}, // 选择的字体
		)
	case "chinese":
		//driver = base64Captcha.NewDriverChinese()
	default:
		driver = base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	}
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
