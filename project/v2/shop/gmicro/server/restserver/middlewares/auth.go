package middlewares

import "github.com/gin-gonic/gin"

// AuthStrategy defines the set of methods used to resource authentication.
// AuthStrategy 定义了用于资源认证的一组方法
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 用于不同认证策略之间的切换
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy used to set to another authentication strategy.
func (ao *AuthOperator) SetStrategy(strategy AuthStrategy) {
	ao.strategy = strategy
}

// AuthFunc execute resource authentication
func (ao *AuthOperator) AuthFunc() gin.HandlerFunc {
	return ao.strategy.AuthFunc()
}
