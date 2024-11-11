package v1

import (
	"context"
	"shop/app/user/srv/data/v1/mock"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"testing"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), metav1.ListMeta{})
}
