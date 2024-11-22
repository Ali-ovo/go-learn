package pkg

import (
	"shop/app/shop/api/internal/data/v1"
	"shop/app/shop/api/internal/service"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/server/restserver/middlewares"
	"shop/pkg/code"
	"shop/pkg/options"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CreateJWT 生成 自定义 token
//
//	@Description:
//	@param user
//	@param opts JWT相关配置
//	@return string JWT token
//	@return error
func CreateJWT(user *data.UserDO, opts *options.JwtOptions) (string, error) {
	// 生成 token
	j := middlewares.NewJWT(opts.Key)

	claims := service.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opts.Realm,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(opts.Timeout)), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Local()),                   // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),                   // 签发时间
		},
	}

	var method jwt.SigningMethod
	switch opts.Method {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	case "ES256":
		method = jwt.SigningMethodES256
	case "ES384":
		method = jwt.SigningMethodES384
	case "ES512":
		method = jwt.SigningMethodES512
	default:
		return "", errors.WithCode(code.ErrJWTDeploy, "invalid jwt method")
	}

	token, err := j.CreateToken(claims, method)
	if err != nil {
		return "", err
	}
	return token, nil
}
