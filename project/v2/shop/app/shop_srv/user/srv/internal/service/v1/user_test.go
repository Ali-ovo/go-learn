package srv

import (
	"context"
	"fmt"
	"shop/app/shop_srv/user/srv/internal/data/v1/mock"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"testing"
)

func TestUserList(t *testing.T) {
	userSrv := NewService(mock.NewMockFactory())
	list, err := userSrv.User().List(context.Background(), []string{}, metav1.ListMeta{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}
