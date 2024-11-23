package controller

import (
	"context"
	user_pb "shop/api/user/v1"
	"shop/gmicro/pkg/errors"
)

// GetUserById
//
//	@Description: 通过 用户id 获取用户信息
//	@receiver uc
//	@param ctx
//	@param request
//	@return *user_pb.UserListResponse
//	@return error
func (uc *userServer) GetUserById(ctx context.Context, request *user_pb.IdRequest) (*user_pb.UserInfoResponse, error) {
	//log.Info("get user by id function called.")
	user, err := uc.srv.GetByID(ctx, uint64(request.Id))
	if err != nil {
		//log.Errorf("get user by id: %s, error: %v", request.Id, err)
		return nil, errors.ToGrpcError(err)
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
