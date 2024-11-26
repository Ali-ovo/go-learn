package controller

import (
	"context"
	user_pb "shop/api/user/v1"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/log"
)

// GetUserByMobile
//
//	@Description: 通过 用户手机号 获取用户信息
//	@receiver uc
//	@param ctx
//	@param request
//	@return *user_pb.UserInfoResponse
//	@return error
func (uc *userServer) GetUserByMobile(ctx context.Context, request *user_pb.MobileRequest) (*user_pb.UserInfoResponse, error) {
	log.Info("get user by mobile function called.")
	user, err := uc.srv.User().GetByMobile(ctx, request.Mobile)
	if err != nil {
		//log.Errorf("get user by mobile: %s, error: %v", request.Mobile, err)
		return nil, errors.ToGrpcError(err)
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
