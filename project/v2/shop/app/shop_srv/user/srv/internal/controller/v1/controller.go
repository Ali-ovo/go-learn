package controller

import (
	user_pb "shop/api/user/v1"
	"shop/app/shop_srv/user/srv/internal/domain/dto"
	"shop/app/shop_srv/user/srv/internal/service"
)

type userServer struct {
	user_pb.UnimplementedUserServer
	srv service.ServiceFactory
}

func NewUserServer(srv service.ServiceFactory) *userServer {
	return &userServer{srv: srv}
}

var _ user_pb.UserServer = &userServer{}

func DTOToResponse(userDTO dto.UserDTO) *user_pb.UserInfoResponse {
	userInfoRsp := user_pb.UserInfoResponse{
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
