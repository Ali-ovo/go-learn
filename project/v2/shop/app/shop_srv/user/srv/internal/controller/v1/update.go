package controller

import (
	"context"
	upbv1 "shop/api/user/v1"
	"shop/pkg/gorm"

	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/service/v1"
	"shop/gmicro/pkg/log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

func (uc *userServer) UpdateUser(ctx context.Context, info *upbv1.UpdateUserInfo) (*empty.Empty, error) {
	//log.Info("update user function called.")

	birthDay := time.Unix(int64(info.BirthDay), 0)

	userDO := data.UserDO{
		BaseModel: gorm.BaseModel{
			ID: info.Id,
		},
		NickName: info.NickName,
		Gender:   info.Gender,
		Birthday: &birthDay,
	}
	userDTO := service.UserDTO{UserDO: userDO}

	err := uc.srv.Update(ctx, &userDTO)
	if err != nil {
		log.Errorf("update user: %v, error: %v", userDTO, err)
		return nil, err
	}
	return &empty.Empty{}, nil
}
