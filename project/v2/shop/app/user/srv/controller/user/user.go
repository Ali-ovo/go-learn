package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	v1 "shop/api/user/v1"
	srvv1 "shop/app/user/srv/service/v1"
)

type userServer struct {
	v1.UnimplementedUserServer
	srv srvv1.UserSrv
}

func (us *userServer) GetUserByMobile(ctx context.Context, request *v1.MobileRequest) (*v1.UserInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userServer) GetUserById(ctx context.Context, request *v1.IdRequest) (*v1.UserInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userServer) CreateUser(ctx context.Context, info *v1.CreateUserInfo) (*v1.UserInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userServer) UpdateUser(ctx context.Context, info *v1.UpdateUserInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userServer) CheckPassWord(ctx context.Context, info *v1.PasswordCheckInfo) (*v1.CheckResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserServer(srv srvv1.UserSrv) *userServer {
	return &userServer{srv: srv}
}
