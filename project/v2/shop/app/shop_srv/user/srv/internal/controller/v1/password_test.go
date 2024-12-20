package controller

import (
	"context"
	"fmt"
	user_pb "shop/api/user/v1"
	"shop/app/shop_srv/user/srv/internal/service"
	"testing"
)

func TestCheckPassWord(t *testing.T) {
	var srv service.ServiceFactory
	userServer := NewUserServer(srv)
	info := &user_pb.PasswordCheckInfo{
		Password:          "56248123",
		EncryptedPassword: "$pbkdf2-sha512$oVWtbs6b1s$5a208122012ce7735ee72ee8b32d7a2b91a648c64e03668246173f53ba558a9e",
	}
	word, err := userServer.CheckPassWord(context.Background(), info)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(word)
}
