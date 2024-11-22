package mock

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type users struct {
	users []*data.UserDO
}

func (u *users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*data.UserDOList, error) {
	users := []*data.UserDO{
		{Mobile: "CZC"},
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

func NewUsers() *users {
	return &users{}
}
