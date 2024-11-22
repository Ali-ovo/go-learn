package controller

import (
	"context"
	"crypto/sha512"
	"fmt"
	upbv1 "shop/api/user/v1"
	dv1 "shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/service/v1"
	"shop/gmicro/pkg/log"

	"github.com/anaskhan96/go-password-encoder"
)

// CreateUser
//
//	@Description: 创建 用户
//	@receiver uc
//	@param ctx
//	@param info
//	@return *upbv1.UserInfoResponse
//	@return error
func (uc *userServer) CreateUser(ctx context.Context, info *upbv1.CreateUserInfo) (*upbv1.UserInfoResponse, error) {
	//log.Info("create user function called.")

	// 密码加密
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(info.Password, options)

	userDO := dv1.UserDO{
		Mobile:   info.Mobile,
		Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd),
		NickName: info.NickName,
	}
	userDTO := service.UserDTO{UserDO: userDO}

	err := uc.srv.Create(ctx, &userDTO)
	if err != nil {
		log.Errorf("create user: %v, error: %v", userDTO, err)
		return nil, err
	}

	userInfoRsp := DTOToResponse(userDTO)

	return userInfoRsp, nil
}
