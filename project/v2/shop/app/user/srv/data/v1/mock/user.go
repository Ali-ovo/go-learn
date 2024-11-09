package mock

import (
	"context"
	dv1 "go-learn/project/v2/shop/app/user/srv/data/v1"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"
)

type users struct {
	users []*dv1.UserDO
}

func NewUsers() *users {
	return &users{}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*dv1.UserDOList, error) {

	return &dv1.UserDOList{
		TotalCount: 1,
		Items:      u.users,
	}, nil
}
