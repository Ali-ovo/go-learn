package BasicAuth

import (
	"shop/app/shop_api/api/service/user/v1"
	"shop/gmicro/server/restserver/middlewares"
	"shop/gmicro/server/restserver/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// NewBasicAuth Basic认证 一般用于 内部认证
//
//	@Description:
//	@param srv
//	@return middlewares.AuthStrategy
func NewBasicAuth(srv user.UserSrv) auth.BasicStrategy {
	return auth.NewBasicStrategy(func(c *gin.Context, username string, password string) bool {
		userDTO, err := srv.MobileLogin(c, username, password)
		if err != nil {
			return false
		}
		// 将登入的用户名 设置进来
		c.Set(middlewares.UsernameKey, userDTO.NickName)
		c.Set(middlewares.KeyUserID, float64(userDTO.ID))
		return true
	})
}
