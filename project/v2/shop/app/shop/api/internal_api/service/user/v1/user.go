package serviceUser

import (
	"context"
	"shop/app/shop/api/internal_api/data/v1"
	doUser "shop/app/shop/api/internal_api/domain/do/user"
	dtoUser "shop/app/shop/api/internal_api/domain/dto/user"
	srvUser "shop/app/shop/api/internal_api/service/user"
	"shop/app/shop/api/pkg/auth/JWTAuth"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/storage"
	"shop/pkg/code"
	"shop/pkg/options"
	"time"
)

type userService struct {
	data    data.DataFactory
	jwtOpts *options.JwtOptions
}

func (us *userService) MobileLogin(ctx context.Context, mobile, password string) (*dtoUser.UserDTO, error) {
	byMobile, err := us.data.User().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	// 检查密码是否正确
	err = us.data.User().CheckPassWord(ctx, password, byMobile.PassWord)
	if err != nil {
		return nil, err
	}

	// 生成 token
	token, err := JWTAuth.CreateJWT(byMobile, us.jwtOpts)
	if err != nil {
		return nil, err
	}

	return &dtoUser.UserDTO{
		UserDO:    *byMobile,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(us.jwtOpts.Timeout).Unix(),
	}, nil
}

func (us *userService) Register(ctx context.Context, mobile, password, codes string) (*dtoUser.UserDTO, error) {
	rstore := storage.RedisCluster{}

	value, err := rstore.GetKey(ctx, mobile)
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不存在")
	}
	if value != codes {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不正确")
	}

	// 注册账号
	var userDO = &doUser.UserDO{
		NickName: mobile,
		Mobile:   mobile,
		PassWord: password,
	}
	err = us.data.User().Create(ctx, userDO)
	if err != nil {
		//log.ErrorfC(ctx, "user register error: %v", err)
		return nil, err
	}

	// 生成 token
	token, err := JWTAuth.CreateJWT(userDO, us.jwtOpts)
	if err != nil {
		return nil, err
	}
	return &dtoUser.UserDTO{
		UserDO:    *userDO,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(us.jwtOpts.Timeout).Unix(),
	}, nil
}

func (us *userService) Update(ctx context.Context, userDTO *dtoUser.UserDTO) error {
	var userDO = &doUser.UserDO{
		ID:       userDTO.ID,
		NickName: userDTO.NickName,
		Birthday: userDTO.Birthday,
		Gender:   userDTO.Gender,
	}
	err := us.data.User().Update(ctx, userDO)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Get(ctx context.Context, userID uint64) (*dtoUser.UserDTO, error) {
	userDO, err := us.data.User().Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dtoUser.UserDTO{UserDO: *userDO}, nil
}

func (us *userService) GetByMobile(ctx context.Context, mobile string) (*dtoUser.UserDTO, error) {
	userDO, err := us.data.User().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &dtoUser.UserDTO{UserDO: *userDO}, nil
}

func NewUserService(data data.DataFactory, jwtOpts *options.JwtOptions) srvUser.UserSrv {
	return &userService{
		data:    data,
		jwtOpts: jwtOpts,
	}
}
