package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"shop/gmicro/server/restserver/middlewares"
)

// AuthzAudience defines the value of jwt audience field.
const AuthzAudience = "shop.czc.com"

// JWTStrategy defines jwt bearer authentication strategy.
// JWTStrategy定义了jwt承载身份验证策略.
type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

var _ middlewares.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy create jwt bearer strategy with GinJWTMiddleware.
func NewJWTStrategy(gjwt ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gjwt}
}

// AuthFunc defines jwt bearer strategy as the gin authentication middleware.
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
