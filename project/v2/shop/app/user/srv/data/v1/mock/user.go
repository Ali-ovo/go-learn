package mock

import (
	"context"
	dv1 "shop/app/user/srv/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type users struct {
	users []*dv1.UserDO
}

func NewUsers() *users {
	return &users{}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	users := []*dv1.UserDO{
		{Name: "CZC"},
	}

	return &dv1.UserDOList{
		TotalCount: 1,
		Items:      users,
	}, nil
}
