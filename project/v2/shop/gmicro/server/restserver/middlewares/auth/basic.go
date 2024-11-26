package auth

import (
	"encoding/base64"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/common/core"
	"strings"

	"shop/gmicro/server/restserver/middlewares"

	"shop/gmicro/pkg/errors"

	"github.com/gin-gonic/gin"
)

// BasicStrategy defines Basic authentication strategy.
type BasicStrategy struct {
	compare func(c *gin.Context, string, password string) bool
}

var _ middlewares.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy create basic strategy with compare function.
func NewBasicStrategy(compare func(c *gin.Context, username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

// AuthFunc defines basic strategy as the gin authentication middleware.
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Authorization 的值 以 空格 切割
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()
			return
		}

		// 虽然是base64编码后发送 但是 base64也是明文传输  一般用于内部访问(因为快捷)
		payload, _ := base64.StdEncoding.DecodeString(auth[1]) // 返回 base64 编码的 二进制数据
		pair := strings.SplitN(string(payload), ":", 2)        // 前面是用户名 后面是密码
		// 核心验证逻辑 在 compare 方法中 验证 账号和密码是否正确  compare 需要注册进来
		if len(pair) != 2 || !b.compare(c, pair[0], pair[1]) {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()
			return
		}
		c.Next()
	}
}
