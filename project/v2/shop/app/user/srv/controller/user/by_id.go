package user

import (
	"context"
	upbv1 "shop/api/user/v1"
)

// GetUserById
//
//	@Description: 通过 用户id 获取用户信息
//	@receiver uc
//	@param ctx
//	@param request
//	@return *upbv1.UserListResponse
//	@return error
func (uc *userServer) GetUserById(ctx context.Context, request *upbv1.IdRequest) (*upbv1.UserInfoResponse, error) {
	//log.Info("get user by id function called.")
	user, err := uc.srv.GetByID(ctx, uint64(request.Id))
	if err != nil {
		//log.Errorf("get user by id: %s, error: %v", request.Id, err)
		return nil, err
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
