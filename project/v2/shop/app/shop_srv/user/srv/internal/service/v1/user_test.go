package service

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/data/v1/mock"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"testing"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), []string{}, metav1.ListMeta{})
}
