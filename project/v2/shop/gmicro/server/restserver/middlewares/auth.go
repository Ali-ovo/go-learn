package middlewares

import "github.com/gin-gonic/gin"

// AuthStrategy defines the set of methods used to resource authentication.
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator used to switch between different authentication strategy.
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
