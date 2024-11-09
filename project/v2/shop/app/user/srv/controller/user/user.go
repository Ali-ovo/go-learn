package user

import (
	srv1 "go-learn/project/v2/shop/app/user/srv/service/v1"
)

type userServer struct {
	srv srv1.UserSrv
}

func NewUserServer(srv srv1.UserSrv) *userServer {
	return &userServer{
		srv: srv,
	}
}
