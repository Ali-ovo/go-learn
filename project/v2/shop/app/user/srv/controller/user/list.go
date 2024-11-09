package user

import (
	"context"
	user_pb "go-learn/project/v2/shop/api/user/v1"
	srvv1 "go-learn/project/v2/shop/app/user/srv/service/v1"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"
)

func DTOToResponse(userdto srvv1.UserDTO) user_pb.UserInfoResponse {
	return user_pb.UserInfoResponse{}
}

func (us *userServer) GetUserList(ctx context.Context, info *user_pb.PageInfo) (*user_pb.UserListResponse, error) {
	srvOpts := metav1.ListMeta{
		Page:     int(info.Pn),
		PageSize: int(info.PSize),
	}

	dtoList, err := us.srv.List(ctx, srvOpts)

	if err != nil {
		return nil, err
	}

	var rsp user_pb.UserListResponse
	for _, value := range dtoList.Items {
		userRsp := DTOToResponse(*value)
		rsp.Data = append(rsp.Data, &userRsp)
	}

	return &rsp, nil
}
