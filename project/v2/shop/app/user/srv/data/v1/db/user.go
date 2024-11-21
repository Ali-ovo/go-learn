package db

import (
	"context"
	"shop/app/user/srv/data/v1"
	code2 "shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"

	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

// List
//
//	@Description: 获取用户列表, 列表页返回 需要返回 Count
//	@receiver u
//	@param ctx
//	@param opts
//	@return *dv1.UserDOList
//	@return error
func (u users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*data.UserDOList, error) {
	// 实现 gorm 查询
	ret := &data.UserDOList{}

	// 处理分页
	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}
	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}

	// 排序
	query := u.db
	for _, v := range orderby {
		query.Order(v)
	}

	// 查询 TODO: 可能存在问题 需要跑一下代码
	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

// GetByMobile
//
//	@Description: 根据手机号获取用户信息
//	@receiver u
//	@param ctx
//	@param mobile: 手机号
//	@return *dv1.UserDO
//	@return error
func (u users) GetByMobile(ctx context.Context, mobile string) (*data.UserDO, error) {
	user := data.UserDO{}

	// err 是 gorm 的error 尽量别往上抛 改成通用的错误 方便后续更换 mysql
	if err := u.db.Where(data.UserDO{Mobile: mobile}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// GetByID
//
//	@Description: 根据 ID 来获取用户信息
//	@receiver u
//	@param ctx
//	@param id: 用户 id
//	@return *dv1.UserDO
//	@return error
func (u users) GetByID(ctx context.Context, id uint64) (*data.UserDO, error) {
	user := data.UserDO{}

	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// Create
//
//	@Description: 创建用户
//	@receiver u
//	@param ctx
//	@param user: 用户 DO 结构体
//	@return error
func (u users) Create(ctx context.Context, user *data.UserDO) error {
	if err := u.db.Create(user).Error; err != nil {
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

// Update
//
//	@Description: 更新用户信息
//	@receiver u
//	@param ctx
//	@param user
//	@return error
func (u users) Update(ctx context.Context, user *data.UserDO) error {
	if err := u.db.Save(user).Error; err != nil {
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

var _ data.UserData = &users{}
