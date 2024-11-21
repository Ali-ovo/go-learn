package mock

import (
	"context"
	"shop/app/user/srv/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type users struct {
	users []*data.UserDO
}

func NewUsers() *users {
	return &users{}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*data.UserDOList, error) {
	users := []*data.UserDO{
		{Mobile: "Ali"},
	}

	return &data.UserDOList{
		TotalCount: 1,
		Items:      users,
	}, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (*data.UserDO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *users) GetByID(ctx context.Context, id uint64) (*data.UserDO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *users) Create(ctx context.Context, user *data.UserDO) error {
	//TODO implement me
	panic("implement me")
}

func (u *users) Update(ctx context.Context, user *data.UserDO) error {
	//TODO implement me
	panic("implement me")
}
