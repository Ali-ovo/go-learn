package service

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserID      uint `json:"userid"`
	NickName    string
	AuthorityId uint
	jwt.RegisteredClaims
}
