package rpc

import (
	"context"
	upbv1 "shop/api/user/v1"
	"shop/app/shop/api/internal/data/v1"
	"shop/gmicro/pkg/common/time"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/registry"
	"shop/gmicro/server/rpcserver"
	"shop/pkg/code"
)

const serviceName = "discovery:///user_srv"

type users struct {
	uc upbv1.UserClient
}

func NewUsers(uc upbv1.UserClient) *users {
	return &users{
		uc: uc,
	}
}

func NewUserServiceClient(r registry.Discovery) upbv1.UserClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		// 设置负载均衡
		rpcserver.WithBanlancerName("selector"),
		// 多添加一个 /  因为 方便做切割 direct:///192.168.16.154:8081 转换成 URL.Path: /192.168.16.154:8081  URL.Scheme: direct
		rpcserver.WithDiscovery(r),
		rpcserver.WithEndpoint(serviceName),
		//rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),		// 这是自己封装的 链路追踪
		rpcserver.WithClientEnableTracing(true),
		//rpc.WithClientTimeout(time.Duration(1000)*time.Second),
	)

	if err != nil {
		panic(err)
	}

	c := upbv1.NewUserClient(conn)
	return c
}

func (u *users) Create(ctx context.Context, user *data.UserDO) error {
	protoUser := &upbv1.CreateUserInfo{
		Mobile:   user.Mobile,
		Password: user.PassWord,
		NickName: user.NickName,
	}
	userRsp, err := u.uc.CreateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	user.ID = int64(userRsp.Id)
	return nil
}

func (u *users) Update(ctx context.Context, user *data.UserDO) error {
	protoUser := &upbv1.UpdateUserInfo{
		Id:       int32(user.ID),
		NickName: user.NickName,
		Gender:   user.Gender,
		BirthDay: uint64(user.Birthday.Unix()),
	}
	_, err := u.uc.UpdateUser(ctx, protoUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *users) Get(ctx context.Context, userID int64) (*data.UserDO, error) {
	protoUser := &upbv1.IdRequest{
		Id: int32(userID),
	}
	user, err := u.uc.GetUserById(ctx, protoUser)
	if err != nil {
		return nil, err
	}
	return &data.UserDO{
		ID:       int64(user.Id),
		NickName: user.NickName,
		Birthday: time.Unix(int64(user.BirthDay), 0),
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
		PassWord: user.Password,
	}, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (*data.UserDO, error) {
	protoUser := &upbv1.MobileRequest{Mobile: mobile}
	user, err := u.uc.GetUserByMobile(ctx, protoUser)
	if err != nil {
		return nil, err
	}
	return &data.UserDO{
		ID:       int64(user.Id),
		NickName: user.NickName,
		Birthday: time.Unix(int64(user.BirthDay), 0),
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
		PassWord: user.Password,
	}, nil
}

func (u *users) CheckPassWord(ctx context.Context, password string, encryptedPwd string) error {
	protoUser := &upbv1.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPwd,
	}
	cres, err := u.uc.CheckPassWord(ctx, protoUser)
	if err != nil {
		return err
	}
	if cres.Success == true {
		return nil
	}
	return errors.WithCode(code.ErrUserPasswordIncorrect, "用户密码错误")
}

var _ data.UserData = &users{}
