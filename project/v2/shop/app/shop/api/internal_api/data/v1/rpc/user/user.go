package user

import (
	"context"
	user_pb "shop/api/user/v1"
	"shop/app/shop/api/internal_api/data/v1/userSrv"
	doUser "shop/app/shop/api/internal_api/domain/do/user"
	"shop/gmicro/pkg/common/time"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"
)

type users struct {
	uc user_pb.UserClient
}

func (u *users) Create(ctx context.Context, user *doUser.UserDO) error {
	protoUser := &user_pb.CreateUserInfo{
		Mobile:   user.Mobile,
		Password: user.PassWord,
		NickName: user.NickName,
	}
	userRsp, err := u.uc.CreateUser(ctx, protoUser)
	if err != nil {
		return errors.FromGrpcError(err)
	}
	user.ID = userRsp.Id
	return nil
}

func (u *users) Update(ctx context.Context, user *doUser.UserDO) error {
	protoUser := &user_pb.UpdateUserInfo{
		Id:       user.ID,
		NickName: user.NickName,
		Gender:   user.Gender,
		BirthDay: uint64(user.Birthday.Unix()),
	}
	_, err := u.uc.UpdateUser(ctx, protoUser)
	if err != nil {
		return errors.FromGrpcError(err)
	}
	return nil
}

func (u *users) Get(ctx context.Context, userID uint64) (*doUser.UserDO, error) {
	protoUser := &user_pb.IdRequest{
		Id: int32(userID),
	}
	user, err := u.uc.GetUserById(ctx, protoUser)
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}
	return &doUser.UserDO{
		ID:       user.Id,
		NickName: user.NickName,
		Birthday: time.Unix(int64(user.BirthDay), 0),
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
		PassWord: user.Password,
	}, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (*doUser.UserDO, error) {
	protoUser := &user_pb.MobileRequest{Mobile: mobile}
	user, err := u.uc.GetUserByMobile(ctx, protoUser)
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}
	return &doUser.UserDO{
		ID:       user.Id,
		NickName: user.NickName,
		Birthday: time.Unix(int64(user.BirthDay), 0),
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
		PassWord: user.Password,
	}, nil
}

func (u *users) CheckPassWord(ctx context.Context, password string, encryptedPwd string) error {
	protoUser := &user_pb.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPwd,
	}
	cres, err := u.uc.CheckPassWord(ctx, protoUser)
	if err != nil {
		return errors.FromGrpcError(err)
	}
	if cres.Success == true {
		return nil
	}
	return errors.WithCode(code2.ErrUserPasswordIncorrect, "用户密码错误")
}

func NewUser(uClient user_pb.UserClient) userSrv.UserData {
	return &users{
		uc: uClient,
	}
}
