package main

import (
	"context"
	"fmt"
	"go-learn/shop/shop_srvs/user_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})

	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)

		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:        "admin123",
			EncryptPassword: user.PassWord,
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("user_%d", i),
			Mobile:   fmt.Sprintf("1350000000%d", i),
			PassWord: "admin123",
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(rsp.Id)
	}

}

func main() {
	Init()

	TestGetUserList()

	// TestCreateUser()

	defer conn.Close()
}
