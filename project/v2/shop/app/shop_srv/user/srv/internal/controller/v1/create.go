package controller

import (
	"context"
	"crypto/sha512"
	"fmt"
	user_pb "shop/api/user/v1"
	"shop/app/shop_srv/user/srv/internal/domain/do"
	"shop/app/shop_srv/user/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"

	"github.com/anaskhan96/go-password-encoder"
)

// CreateUser
//
//	@Description: 创建 用户
//	@receiver uc
//	@param ctx
//	@param info
//	@return *user_pb.UserInfoResponse
//	@return error
func (uc *userServer) CreateUser(ctx context.Context, info *user_pb.CreateUserInfo) (*user_pb.UserInfoResponse, error) {
	log.Info("create user function called.")

	// 密码加密
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(info.Password, options)

	userDO := do.UserDO{
		Mobile:   info.Mobile,
		Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd),
		NickName: info.NickName,
	}
	userDTO := dto.UserDTO{UserDO: userDO}

	err := uc.srv.User().Create(ctx, &userDTO)
	if err != nil {
		//log.Errorf("create user: %v, error: %v", userDTO, err)
		return nil, errors.ToGrpcError(err)
	}

	userInfoRsp := DTOToResponse(userDTO)

	return userInfoRsp, nil
}
