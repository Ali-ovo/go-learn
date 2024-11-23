package controller

import (
	"context"
	"crypto/sha512"
	user_pb "shop/api/user/v1"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

func (uc *userServer) CheckPassWord(ctx context.Context, info *user_pb.PasswordCheckInfo) (*user_pb.CheckResponse, error) {
	//校验密码
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(info.EncryptedPassword, "$")
	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
	return &user_pb.CheckResponse{Success: check}, nil
}
