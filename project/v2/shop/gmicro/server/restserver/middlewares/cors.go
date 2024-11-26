package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		// 允许跨域访问的域名：若有端口需写全（协议+域名+端口）若没有端口末尾不用加'/' 用通配符*表示允许任何域的JavaScript访问资源
		ctx.Header("Access-Control-Allow-Origin", "*")
		// 提示OPTIONS预检时，后端需要设置的两个常用自定义头
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, x-token")
		// 指定对预请求的响应中，哪些 HTTP 方法允许访问请求的资源。 * 代办允许所有方法
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		// 指示哪些 HTTP 头的名称能在响应中列出
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 表示是否允许发送Cookie
		ctx.Header("Access-Control-Allow-Credentials", "true") // 表示是否允许发送Cookie

		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
