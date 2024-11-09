package v1

import (
	"context"
	dv1 "go-learn/project/v2/shop/app/user/srv/data/v1"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"
)

type UserDTO struct {
	dv1.UserDO
}

type UserSrv interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error)
}

type userService struct {
	userStore dv1.UserStore
}

func NewUserService(us dv1.UserStore) *userService {
	return &userService{userStore: us}
}

var _ UserSrv = &userService{}

type UserDTOList struct {
	TotalCount int64      `json:"totalCount,omitempty"`
	Items      []*UserDTO `json:"data"`
}

func (u *userService) List(ctx context.Context, opts metav1.ListMeta) (*UserDTOList, error) {
	doList, err := u.userStore.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	var userDTOList UserDTOList

	for _, value := range doList.Items {
		userDTO := UserDTO{*value}
		userDTOList.Items = append(userDTOList.Items, &userDTO)
	}

	return &userDTOList, nil
}
