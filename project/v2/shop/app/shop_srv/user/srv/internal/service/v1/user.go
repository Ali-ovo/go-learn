package srv

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/domain/dto"
	"shop/app/shop_srv/user/srv/internal/service"
	code2 "shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"

	"gorm.io/gorm"
)

type userService struct {
	data data.DataFactory
}

func (u *userService) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*dto.UserDTOList, error) {
	// 业务逻辑1

	/*
		1. data层的接口必须先写好
		2. 在写测试案例的时候每次测试底层的data层的数据按照我期望的返回
			1. 先手动去插入一些数据
			2. 去删除一些数据
		3. 如果 data 层的方法有bug 代码想要具备 良好可测试性
	*/
	doList, err := u.data.User().List(ctx, orderby, opts)
	if err != nil {
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	// 业务逻辑2
	// 代码不方便 会导致写单元测试用例难写
	var userDTOList dto.UserDTOList
	for _, value := range doList.Items {
		projectDTO := dto.UserDTO{UserDO: *value}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	return &userDTOList, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*dto.UserDTO, error) {
	userDo, err := u.data.User().GetByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return &dto.UserDTO{UserDO: *userDo}, nil
}

func (u *userService) GetByID(ctx context.Context, id uint64) (*dto.UserDTO, error) {
	userDo, err := u.data.User().GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return &dto.UserDTO{UserDO: *userDo}, nil
}

func (u *userService) Create(ctx context.Context, user *dto.UserDTO) error {
	// 先判断用户号码是否存在
	_, err := u.data.User().GetByMobile(ctx, user.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if result := u.data.User().Create(ctx, nil, &user.UserDO); result != nil {
				if result.Error != nil {
					return errors.WithCode(code2.ErrDatabase, err.Error())
				}
				return errors.WithCode(code.ErrUserNotFound, "Create User failure")
			}
			return nil
		}
		// 说明 数据库存在问题
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return errors.WithCode(code.ErrUserAlreadyExists, "用户已经存在")
}

func (u *userService) Update(ctx context.Context, user *dto.UserDTO) error {
	// 先判断用户id 是否存在
	userDO, err := u.data.User().GetByID(ctx, uint64(user.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrUserAlreadyExists, "用户不存在")
		}
		// 说明 数据库存在问题
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}

	userDO.NickName = user.NickName
	userDO.Birthday = user.Birthday
	userDO.Gender = user.Gender

	if result := u.data.User().Update(ctx, nil, userDO); result.RowsAffected == 0 {
		if result.Error != nil {
			return errors.WithCode(code2.ErrDatabase, result.Error.Error())
		}
		return errors.WithCode(code.ErrUserNotFound, "用户不存在")
	}
	return nil
}

func newUser(srv *serviceFactory) service.UserSrv {
	return &userService{
		data: srv.data,
	}
}
