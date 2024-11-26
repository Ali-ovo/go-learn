package mock

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type mockFactory struct{}

func (mf *mockFactory) User() data.UserStore {
	return newUsers()
}

func (mf *mockFactory) Begin() *gorm.DB {
	panic("implement me")
}

func NewMockFactory() *mockFactory {
	return &mockFactory{}
}

type users struct {
	users []*do.UserDO
}

func (u *users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.UserDOList, error) {
	users := []*do.UserDO{
		{Mobile: "CZC"},
	}

	return &do.UserDOList{
		TotalCount: 1,
		Items:      users,
	}, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (*do.UserDO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *users) GetByID(ctx context.Context, id uint64) (*do.UserDO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *users) Create(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func (u *users) Update(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func newUsers() data.UserStore {
	return &users{}
}
