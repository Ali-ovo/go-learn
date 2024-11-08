// Package auth encrypt and compare password string.
package auth

import (
	"github.com/dgrijalva/jwt-go/v4"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Sign issue a jwt token based on secretID, secretKey, iss and aud.
// 签发一个 时间较短的 jwt token
func Sign(secretID string, secretKey string, iss, aud string) string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(), // 多久过期
		"iat": time.Now().Unix(),                  // 签发时间
		"nbf": time.Now().Add(0).Unix(),           // 生效时间
		"aud": aud,                                // 签发人
		"iss": iss,                                // 接收人
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = secretID

	// Sign the token with the specified secret.
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}
