package db

import (
	"context"
	dv1 "shop/app/user/srv/data/v1"
	code2 "shop/gmicro/pkg/code"
	v1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	"shop/pkg/code"

	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

// Create implements v1.UserStore.
func (u *users) Create(ctx context.Context, user *dv1.UserDO) error {
	if err := u.db.Create(user).Error; err != nil {
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

// GetByID implements v1.UserStore.
func (u *users) GetByID(ctx context.Context, id uint64) (*dv1.UserDO, error) {
	user := dv1.UserDO{}

	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// GetByMobile implements v1.UserStore.
func (u *users) GetByMobile(ctx context.Context, mobile string) (*dv1.UserDO, error) {
	user := dv1.UserDO{}
	// err 是 gorm 的error 尽量别往上抛 改成通用的错误 方便后续更换 mysql
	if err := u.db.Where(dv1.UserDO{Mobile: mobile}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// List implements v1.UserStore.
func (u *users) List(ctx context.Context, orderby []string, opts v1.ListMeta) (*dv1.UserDOList, error) {
	// 实现 gorm 查询
	ret := &dv1.UserDOList{}

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

// Update implements v1.UserStore.
func (u *users) Update(ctx context.Context, user *dv1.UserDO) error {
	if err := u.db.Save(user).Error; err != nil {
		return errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return nil
}

func NewUsers(db *gorm.DB) *users {
	return &users{db: db}
}

var _ dv1.UserStore = &users{}
