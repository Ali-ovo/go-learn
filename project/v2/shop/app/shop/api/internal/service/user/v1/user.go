package user

import (
	"context"
	"shop/app/shop/api/internal/data/v1"

	"errors"
	"shop/app/shop/api/internal/service"
	"shop/gmicro/server/restserver/middlewares"
	"shop/pkg/options"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	Register(ctx context.Context, mobile, password string) (*UserDTO, error)
	// Update 更新 用户信息
	Update(ctx context.Context, userDTO *UserDTO) error
	// Get 通过 ID 获取 用户信息
	Get(ctx context.Context, userID uint64) (*UserDTO, error)
	// GetByMobile 通过 手机号 获取 用户信息
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	// CheckPassWord 验证 密码是否正确
	CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error)
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
	j := middlewares.NewJWT(us.jwtOpts.Key)

	claims := service.CustomClaims{
		ID:          uint(byMobile.ID),
		NickName:    byMobile.NickName,
		AuthorityId: uint(byMobile.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    us.jwtOpts.Realm,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(us.jwtOpts.Timeout)), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Local()),                         // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),                         // 签发时间
		},
	}

	var method jwt.SigningMethod
	switch us.jwtOpts.Method {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	case "HS512":
		method = jwt.SigningMethodHS512
	case "ES256":
		method = jwt.SigningMethodES256
	case "ES384":
		method = jwt.SigningMethodES384
	case "ES512":
		method = jwt.SigningMethodES512
	default:
		return nil, errors.New("invalid jwt method")
	}

	token, err := j.CreateToken(claims, method)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		UserDO:    *byMobile,
		Token:     token,
		ExpiredAt: time.Now().Local().Add(us.jwtOpts.Timeout).Unix(),
	}, nil
}

func (us *userService) Register(ctx context.Context, mobile, password string) (*UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userService) Update(ctx context.Context, userDTO *UserDTO) error {
	//TODO implement me
	panic("implement me")
}

func (us *userService) Get(ctx context.Context, userID uint64) (*UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (us *userService) CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ UserSrv = (*userService)(nil)
