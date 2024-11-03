package mock

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"go-learn/ch19/mock"
)

func TestGetUserByMobile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserData := NewMockUserData(ctrl)
	mockUserData.EXPECT().GetUserByMobile(gomock.Any(), "18").Return(mock.User{
		NickName: "ali_18",
	}, nil)

	// 实际调用
	userServer := mock.UserServer{
		DB: mockUserData,
	}
	user, err := userServer.GetUserByMobile(context.Background(), "18")

	// 判断正确与否
	if err != nil {
		t.Fatalf("GetUserByMobile error: %v", err)
	}

	if user.NickName == "ali_18" {
		t.Fatalf("GetUserByMobile error: %v", user)
	}

	t.Logf("user: %v", user)
}

func TestGetUserByMobileFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserData := NewMockUserData(ctrl)
	mockUserData.EXPECT().GetUserByMobile(gomock.Any(), "19").Return(mock.User{
		NickName: "ali_19",
	}, nil)

	// 实际调用
	userServer := mock.UserServer{
		DB: mockUserData,
	}
	user, err := userServer.GetUserByMobile(context.Background(), "19")

	// 判断正确与否
	if err != nil {
		t.Fatalf("GetUserByMobile error: %v", err)
	}

	if user.NickName != "ali_19" {
		t.Fatalf("GetUserByMobile error: %v", user)
	}

	t.Logf("user: %v", user)
}
