package controller

import (
	"context"
	user_pb "shop/api/user/v1"
	"shop/app/shop_srv/user/srv/internal/domain/do"
	"shop/app/shop_srv/user/srv/internal/domain/dto"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
	"shop/pkg/gorm"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

func (uc *userServer) UpdateUser(ctx context.Context, info *user_pb.UpdateUserInfo) (*empty.Empty, error) {
	log.Info("update user function called.")

	birthDay := time.Unix(int64(info.BirthDay), 0)

	userDO := do.UserDO{
		BaseModel: gorm.BaseModel{
			ID: info.Id,
		},
		NickName: info.NickName,
		Gender:   info.Gender,
		Birthday: &birthDay,
	}
	userDTO := dto.UserDTO{UserDO: userDO}

	err := uc.srv.User().Update(ctx, &userDTO)
	if err != nil {
		//log.Errorf("update user: %v, error: %v", userDTO, err)
		return nil, errors.ToGrpcError(err)
	}
	return &empty.Empty{}, nil
}
