package user

import (
	upbv1 "shop/api/user/v1"
	srvv1 "shop/app/user/srv/service/v1"
)

type userServer struct {
	upbv1.UnimplementedUserServer
	srv srvv1.UserSrv
}

func NewUserServer(srv srvv1.UserSrv) *userServer {
	return &userServer{srv: srv}
}

func DTOToResponse(userDTO srvv1.UserDTO) *upbv1.UserInfoResponse {
	// 在 grpc 的 message 中字段有默认值, 你不能随便赋值 nil 进去, 容易出错
	// 这里要搞清, 那些字段是有默认值
	userInfoRsp := upbv1.UserInfoResponse{
		Id:       userDTO.ID,
		Mobile:   userDTO.Mobile,
		Password: userDTO.Password,
		NickName: userDTO.NickName,
		Gender:   userDTO.Gender,
		Role:     uint32(userDTO.Role),
	}
	if userDTO.Birthday != nil {
		userInfoRsp.BirthDay = uint64(userDTO.Birthday.Unix())
	}
	return &userInfoRsp
}

var _ upbv1.UserServer = &userServer{}
