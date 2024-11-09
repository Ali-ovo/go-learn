package db

import (
	"context"
	dv1 "go-learn/project/v2/shop/app/user/srv/data/v1"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

func newUsers(db *gorm.DB) *users {
	return &users{db: db}
}

func (u *users) List(ctx context.Context, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	return nil, nil
}
