package auth

import (
	"errors"
	"fmt"
	"go-learn/project/v2/shop/pkg/common/util/idutil"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func TestSign(t *testing.T) {
	secretID := idutil.NewSecretID()
	SigningKey := "123456"
	tokenString := Sign(secretID, SigningKey, "czc", "zsz")

	// 解析 token 中的数据
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// 因为错误类型是 uInt32 类型 并且 这些类型是 2进制的数字:1,2,4,8,16... 所以使用 & 比较更快捷
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				panic(TokenMalformed)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				panic(TokenExpired)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				panic(TokenNotValidYet)
			} else {
				panic(TokenInvalid)
			}
		}
	}
	fmt.Println(token)

}
