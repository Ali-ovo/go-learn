package middlewares

import "github.com/gin-gonic/gin"

const (
	UsernameKey = "username"
	KeyUserID   = "userid"
	UserIP      = "ip"
)

// Context 为每个请求添加上下文, 类型 django 中间件
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 类似一下操作 获取我们想要的数据 可以在所有的业务逻辑中 c.Get 拿取数据
		// 也可以写到 jwt 中  然后 获取 相关 ip
		// 从 c 中获取到 ip 地址
		c.Set(UserIP, c.ClientIP())
		c.Next()
	}
}
