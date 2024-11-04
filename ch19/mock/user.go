package mock

import "context"

type User struct {
	Mobile   string
	Password string
	NickName string
}

type UserServer struct {
	DB UserData
}

type UserData interface {
	GetUserByMobile(ctx context.Context, mobile string) (User, error)
}

func (us *UserServer) GetUserByMobile(ctx context.Context, mobile string) (User, error) {
	user, err := us.DB.GetUserByMobile(ctx, mobile)
	if err != nil {
		return User{}, err
	}

	if user.NickName == "ali_18" {
		user.NickName = "ali_17"
	}

	return user, nil
}
