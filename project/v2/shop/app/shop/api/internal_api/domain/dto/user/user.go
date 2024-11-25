package dtoUser

import (
	doUser "shop/app/shop/api/internal_api/domain/do/user"
)

type UserDTO struct {
	doUser.UserDO
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}
