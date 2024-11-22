package controller

import (
	"context"
	upbv1 "shop/api/user/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
)

/*
controller 层依赖了 service, service 层依赖 data 层
	controller 层 可以直接依赖 data 层
controller 依赖了 service 并不是直接依赖了具体的 struct 而是依赖了 interface 好处 随时更换  做到解耦

// java 中的 ioc, 控制翻转 ioc = invertion of control
// 代码分层, 第三方服务, rpc, redis, 等等, 带来一定的复杂度
*/

func (uc *userServer) GetUserList(ctx context.Context, request *upbv1.PageInfo) (*upbv1.UserListResponse, error) {
	//log.Info("GetUserList is called")
	srvOpts := metav1.ListMeta{
		Page:     int(request.Pn),
		PageSize: int(request.PSize),
	}
	dtoList, err := uc.srv.List(ctx, []string{}, srvOpts)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var rsp upbv1.UserListResponse
	for _, value := range dtoList.Items {
		userRsp := DTOToResponse(*value)
		rsp.Data = append(rsp.Data, userRsp)
	}

	return &rsp, nil
}
