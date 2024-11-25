package srvUser

import (
	"context"
	dtoUser "shop/app/shop_api/api/internal/domain/dto/user"
)

type UserSrv interface {
	// MobileLogin 密码登入
	MobileLogin(ctx context.Context, mobile, password string) (*dtoUser.UserDTO, error)
	// Register 注册用户账号
	Register(ctx context.Context, mobile, password, code string) (*dtoUser.UserDTO, error)
	// Update 更新 用户信息
	Update(ctx context.Context, userDTO *dtoUser.UserDTO) error
	// Get 通过 ID 获取 用户信息
	Get(ctx context.Context, userID uint64) (*dtoUser.UserDTO, error)
	// GetByMobile 通过 手机号 获取 用户信息
	GetByMobile(ctx context.Context, mobile string) (*dtoUser.UserDTO, error)
}
