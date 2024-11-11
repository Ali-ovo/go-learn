package auth

import (
	"fmt"
	"shop/gmicro/pkg/code"
	"time"

	"shop/gmicro/pkg/common/core"
	"shop/gmicro/server/restserver/middlewares"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"shop/gmicro/pkg/errors"
)

// Defined errors.
var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

// Secret contains the basic information of the secret key.
// Secret 包含 密钥的基本信息
type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

// CacheStrategy defines jwt bearer authentication strategy which called `cache strategy`.
// Secrets are obtained through grpc api interface and cached in memory.
type CacheStrategy struct {
	get func(kid string) (Secret, error)
}

var _ middlewares.AuthStrategy = &CacheStrategy{}

// NewCacheStrategy create cache strategy with function which can list and cache secrets.
func NewCacheStrategy(get func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{get}
}

// AuthFunc defines cache strategy as the gin authentication middleware.
func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Authorization 的值 以 空格 切割
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			core.WriteResponse(c, errors.WithCode(code.ErrMissingHeader, "Authorization header cannot be empty."), nil)
			c.Abort()

			return
		}

		var rawJWT string
		// Parse the header to get the token part.
		// 解析 header 将数据填充到 rawJWT
		fmt.Sscanf(header, "Bearer %s", &rawJWT)

		// Use own validation logic, see below
		// JWT 携带的信息
		var secret Secret
		// 自定义 签名内容 (只需要在 jwt.MapClaims{} 上 在 封装一层)
		claims := &jwt.MapClaims{}
		// Verify the token
		// 验证 token
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is HMAC signature
			// 判断 token 签署方法 是对称加密 HS
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// 返回 token 塞入一个 kid  这里获取 kid
			kid, ok := token.Header["kid"].(string)
			if !ok {
				// 没有 kid
				return nil, ErrMissingKID
			}

			//我们的jwt的以前的认证方式是， 只要解密成功，就认为是合法的
			//如果我有个恶意的用户，他可以伪造一个jwt，然后把kid设置成一个不存在的kid，这样就可以绕过认证，我们可以在token中放字符串
			//我们想要拉黑一个用户
			var err error
			// 如果存在 kid 去 cache 拿去  cache 可以是 redis 等 只要能获取到即可
			secret, err = cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}
			// 如果 key 值
			return []byte(secret.Key), nil
		}, jwt.WithAudience(AuthzAudience)) // 用于设置 JWT 中的 "aud" (audience) 声明(解析的时候也需要携带 不携带 解析不出来)
		if err != nil || !parsedT.Valid {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()

			return
		}

		// 设置过期
		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			core.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Set(middlewares.UsernameKey, secret.Username)
		c.Next()
	}
}

// KeyExpired checks if a key has expired, if the value of user.SessionState.Expires is 0, it will be ignored.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0)) // 判断是否过期
	}

	return false
}
