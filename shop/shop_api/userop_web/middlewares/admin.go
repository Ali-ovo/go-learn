package middlewares

import (
	"go-learn/shop/shop_api/userop_web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdminAuth() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}
