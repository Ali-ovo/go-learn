package v1

import (
	"context"
	dv1 "shop/app/user/srv/data/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"
)

type UserDTO struct {
	// 这里偷个懒, 应为业务层和 底层 字段没有太大变动
	dv1.UserDO
}

type UserSrv interface {
	// List 获取 用户列表页
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error)
	// GetByMobile 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	// GetByID 通过 id 查询用户
	GetByID(ctx context.Context, id uint64) (*UserDTO, error)
	// Create 添加用户
	Create(ctx context.Context, user *UserDTO) error
	// Update 更新用户
	Update(ctx context.Context, user *UserDTO) error
}

type userService struct {
	userStrore dv1.UserStore // 数据的来源
}

var _ UserSrv = &userService{}

// UserDTOList 返回 自定义的结构体 解耦
type UserDTOList struct {
	TotalCount int64      `json:"totalCount,omitempty"` // 总数
	Items      []*UserDTO `json:"data"`                 // 数据
}

func (u *userService) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error) {
	// 业务逻辑1

	/*
		1. data层的接口必须先写好
		2. 在写测试案例的时候每次测试底层的data层的数据按照我期望的返回
			1. 先手动去插入一些数据
			2. 去删除一些数据
		3. 如果 data 层的方法有bug 代码想要具备好的可测试性
	*/
	doList, err := u.userStrore.List(ctx, orderby, opts)
	if err != nil {
		return nil, err
	}
	// 业务逻辑2
	// 代码不方便写单元测试用例
	// 代码不方便 会导致写单元测试用例难写
	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{*value}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	// 业务逻辑3

	return &userDTOList, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	userDo, err := u.userStrore.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	return &UserDTO{UserDO: *userDo}, nil
}

func (u *userService) GetByID(ctx context.Context, id uint64) (*UserDTO, error) {
	userDo, err := u.userStrore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserDTO{UserDO: *userDo}, nil
}

func (u *userService) Create(ctx context.Context, user *UserDTO) error {
	// 先判断用户号码是否存在
	_, err := u.userStrore.GetByMobile(ctx, user.Mobile)
	if errors.IsCode(err, code.ErrUserNotFound) {
		return u.userStrore.Create(ctx, &user.UserDO)
	}

	// 说明 数据库存在问题
	if err != nil {
		return err
	}

	return errors.WithCode(code.ErrUserAlreadyExists, "用户已经存在")
}

func (u *userService) Update(ctx context.Context, user *UserDTO) error {
	// 先判断用户id 是否存在
	_, err := u.userStrore.GetByID(ctx, uint64(user.ID))
	if errors.IsCode(err, code.ErrUserNotFound) {
		return errors.WithCode(code.ErrUserAlreadyExists, "用户不存在")
	}
	// 说明 数据库存在问题
	if err != nil {
		return err
	}

	return u.userStrore.Update(ctx, &user.UserDO)
}

func NewUserService(us dv1.UserStore) *userService {
	return &userService{
		userStrore: us,
	}
}
