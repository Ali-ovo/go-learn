package user

import (
	"context"
	"shop/app/shop_api/api/data/v1"
	"shop/app/shop_api/api/pkg/auth/JWTAuth"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/pkg/storage"
	"shop/pkg/code"
	"shop/pkg/options"
	"time"
)

type UserDTO struct {
	data.UserDO
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}

type UserSrv interface {
	// MobileLogin 密码登入
	MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error)
	// Register 注册用户账号
	Register(ctx context.Context, mobile, password, code string) (*UserDTO, error)
	// Update 更新 用户信息
	Update(ctx context.Context, userDTO *UserDTO) error
	// Get 通过 ID 获取 用户信息
	Get(ctx context.Context, userID uint64) (*UserDTO, error)
	// GetByMobile 通过 手机号 获取 用户信息
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
}

type userService struct {
	data.UserData
	jwtOpts *options.JwtOptions
}

func NewUserService(ud data.UserData, jwtOpts *options.JwtOptions) UserSrv {
	return &userService{
		UserData: ud,
		jwtOpts:  jwtOpts,
	}
}

func (us *userService) MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error) {
	byMobile, err := us.UserData.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	// 检查密码是否正确
	err = us.UserData.CheckPassWord(ctx, password, byMobile.PassWord)
	if err != nil {
		return nil, err
	}

	// 生成 token
	token, err := JWTAuth.CreateJWT(byMobile, us.jwtOpts)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		UserDO:    *byMobile,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(us.jwtOpts.Timeout).Unix(),
	}, nil
}

func (us *userService) Register(ctx context.Context, mobile, password, codes string) (*UserDTO, error) {
	rstore := storage.RedisCluster{}

	value, err := rstore.GetKey(ctx, mobile)
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不存在")
	}
	if value != codes {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不正确")
	}

	// 注册账号
	var userDO = &data.UserDO{
		NickName: mobile,
		Mobile:   mobile,
		PassWord: password,
	}
	err = us.UserData.Create(ctx, userDO)
	if err != nil {
		// log.ErrorfC(ctx, "user register error: %v", err)
		return nil, err
	}

	// 生成 token
	token, err := JWTAuth.CreateJWT(userDO, us.jwtOpts)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		UserDO:    *userDO,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(us.jwtOpts.Timeout).Unix(),
	}, nil
}

func (us *userService) Update(ctx context.Context, userDTO *UserDTO) error {
	var userDO = &data.UserDO{
		ID:       userDTO.ID,
		NickName: userDTO.NickName,
		Birthday: userDTO.Birthday,
		Gender:   userDTO.Gender,
	}
	err := us.UserData.Update(ctx, userDO)
	if err != nil {
		return err
	}
	return nil
}

func (us *userService) Get(ctx context.Context, userID uint64) (*UserDTO, error) {
	userDO, err := us.UserData.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserDTO{UserDO: *userDO}, nil
}

func (us *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	userDO, err := us.UserData.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &UserDTO{UserDO: *userDO}, nil
}

var _ UserSrv = (*userService)(nil)
