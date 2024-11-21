package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth(signKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		j := NewJWT(signKey)
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				if err == TokenExpired {
					c.JSON(http.StatusUnauthorized, map[string]string{
						"msg": "授权已过期",
					})
					c.Abort()
					return
				}
			}

			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(signKey string) *JWT {
	return &JWT{
		[]byte(signKey), // 可以设置过期时间
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims any, method jwt.SigningMethod) (string, error) {
	mapClaims, ok := claims.(*jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	token := jwt.NewWithClaims(method, mapClaims) // 这里选择加密方式 我选择 HS512	 claims 设置需要携带可被查看的信息 和 jwt规定的信息
	return token.SignedString(j.SigningKey)       // 返回 jwt
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	// 是否 jwt值是否有效
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// 因为错误类型是 uInt32 类型 并且 这些类型是 2进制的数字:1,2,4,8,16... 所以使用 & 比较更快捷
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		// 需要额外判断 token.Valid 是否 和 密钥匹配成功
		// 因为就算没有匹配成功  基础数据也是能被解析到的
		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string, method jwt.SigningMethod) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	// 验证 jwt值是否 有效
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
		return j.CreateToken(*claims, method)
	}
	return "", TokenInvalid
}
