package user

import (
	"context"
	upbv1 "shop/api/user/v1"
)

// GetUserByMobile
//
//	@Description: 通过 用户手机号 获取用户信息
//	@receiver uc
//	@param ctx
//	@param request
//	@return *upbv1.UserInfoResponse
//	@return error
func (uc *userServer) GetUserByMobile(ctx context.Context, request *upbv1.MobileRequest) (*upbv1.UserInfoResponse, error) {
	//log.Info("get user by mobile function called.")
	user, err := uc.srv.GetByMobile(ctx, request.Mobile)
	if err != nil {
		//log.Errorf("get user by mobile: %s, error: %v", request.Mobile, err)
		return nil, err
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
