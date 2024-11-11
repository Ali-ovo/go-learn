package db

import (
	"context"
	"gorm.io/gorm"
	dv1 "shop/app/user/srv/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	// 实现 gorm 查询
	return &dv1.UserDOList{}, nil
}
