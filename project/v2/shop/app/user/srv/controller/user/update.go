package user

import (
	"context"
	upbv1 "shop/api/user/v1"
	dv1 "shop/app/user/srv/data/v1"
	srvv1 "shop/app/user/srv/service/v1"
	"shop/gmicro/pkg/log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

func (uc *userServer) UpdateUser(ctx context.Context, info *upbv1.UpdateUserInfo) (*empty.Empty, error) {
	//log.Info("update user function called.")

	birthDay := time.Unix(int64(info.BirthDay), 0)

	userDO := dv1.UserDO{
		BaseModel: dv1.BaseModel{
			ID: info.Id,
		},
		NickName: info.NickName,
		Gender:   info.Gender,
		Birthday: &birthDay,
	}
	userDTO := srvv1.UserDTO{UserDO: userDO}

	err := uc.srv.Update(ctx, &userDTO)
	if err != nil {
		log.Errorf("update user: %v, error: %v", userDTO, err)
		return nil, err
	}
	return &empty.Empty{}, nil
}
