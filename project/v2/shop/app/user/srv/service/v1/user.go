package v1

import (
	"context"
	dv1 "shop/app/user/srv/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type UserSrv interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error)
}

type userService struct {
	userStrore dv1.UserStore // 数据的来源
}

func NewUserService(us dv1.UserStore) *userService {
	return &userService{
		userStrore: us,
	}
}

type UserDTO struct {
	Name string `json:"name"`
}

// UserDTOList 返回 自定义的结构体 解耦
type UserDTOList struct {
	TotalCount int64      `json:"totalCount,omitempty"` // 总数
	Items      []*UserDTO `json:"data"`                 // 数据
}

var _ UserSrv = &userService{}

func (u *userService) List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error) {
	// 业务逻辑1

	/*
		1. data层的接口必须先写好
		2. 我期望测试的时候每次测试底层的data层的数据按照我期望的返回
			1. 先手动去插入一些数据
			2. 去删除一些数据
		3. 如果 data 层的方法有bug 代码想要具备好的可测试性
	*/
	doList, err := u.userStrore.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	// 业务逻辑2
	// 代码不方便写单元测试用例
	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{
			Name: value.Name,
		}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	// 业务逻辑3

	return &userDTOList, nil
}
