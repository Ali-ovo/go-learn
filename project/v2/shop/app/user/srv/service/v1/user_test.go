package v1

import (
	"context"
	"go-learn/project/v2/shop/app/user/srv/data/v1/mock"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"
	"testing"
)

func TestUserList(t *testing.T) {
	userSrv := NewUserService(mock.NewUsers())
	userSrv.List(context.Background(), metav1.ListMeta{})
}
