package userSrv

import (
	"context"
	doUser "shop/app/shop_api/api/internal/domain/do/user"
)

type UserData interface {
	Create(ctx context.Context, user *doUser.UserDO) error
	Update(ctx context.Context, user *doUser.UserDO) error
	Get(ctx context.Context, userID uint64) (*doUser.UserDO, error)
	GetByMobile(ctx context.Context, mobile string) (*doUser.UserDO, error)
	CheckPassWord(ctx context.Context, password string, encryptedPwd string) error
}
