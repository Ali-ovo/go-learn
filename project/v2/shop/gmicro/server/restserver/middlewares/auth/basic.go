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

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

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
