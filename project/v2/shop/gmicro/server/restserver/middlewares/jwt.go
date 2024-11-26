package middlewares

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtOption func(j *JWT)

type JWT struct {
	SigningAlgorithm string
	HsSigningKey     []byte
	EsPrivKey        *ecdsa.PrivateKey
	RsPrivKey        *rsa.PrivateKey
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT(opts ...JwtOption) *JWT {
	jwt := &JWT{}

	for _, o := range opts {
		o(jwt)
	}

	return jwt
}

func WithHsSymmetricEncrypt(alg string, key []byte) JwtOption {
	return func(j *JWT) {
		j.SigningAlgorithm = alg
		j.HsSigningKey = key
	}
}

func WithEsAsymmetricEncrypt(alg string, privkey *ecdsa.PrivateKey) JwtOption {
	return func(j *JWT) {
		j.SigningAlgorithm = alg
		j.EsPrivKey = privkey
	}
}

func WithRsAsymmetricEncrypt(alg string, privkey *rsa.PrivateKey) JwtOption {
	return func(j *JWT) {
		j.SigningAlgorithm = alg
		j.RsPrivKey = privkey
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims any) (string, error) {
	mapClaims, ok := claims.(jwt.Claims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.SigningAlgorithm), mapClaims) // 这里选择加密方式 我选择 HS512	 claims 设置需要携带可被查看的信息 和 jwt规定的信息
	if j.HsSigningKey != nil {
		return token.SignedString(j.HsSigningKey) // 返回 Hs jwt
	} else if j.EsPrivKey != nil {
		return token.SignedString(j.EsPrivKey) // 返回 Es jwt
	} else {
		return token.SignedString(j.RsPrivKey) // 返回 RS jwt
	}
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	// 是否 jwt值是否有效
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.HsSigningKey, nil
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
		return j.HsSigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
