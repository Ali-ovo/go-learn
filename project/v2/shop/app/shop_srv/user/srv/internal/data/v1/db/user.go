package db

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

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
func (u *users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.UserDOList, error) {
	db := u.db.WithContext(ctx)
	var ret do.UserDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.UserDOList{})

	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result = result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count
	return &ret, nil
}

// GetByMobile
//
//	@Description: 根据手机号获取用户信息
//	@receiver u
//	@param ctx
//	@param mobile: 手机号
//	@return *dv1.UserDO
//	@return error
func (u *users) GetByMobile(ctx context.Context, mobile string) (*do.UserDO, error) {
	db := u.db.WithContext(ctx)
	var user do.UserDO

	// err 是 gorm 的error 尽量别往上抛 改成通用的错误 方便后续更换 mysql
	if err := db.Where(do.UserDO{Mobile: mobile}).First(&user).Error; err != nil {
		return nil, err
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
func (u *users) GetByID(ctx context.Context, id uint64) (*do.UserDO, error) {
	db := u.db.WithContext(ctx)
	var user do.UserDO

	if err := db.First(&user, id).Error; err != nil {
		return nil, err
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
func (u *users) Create(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB {
	db := u.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(user)
}

// Update
//
//	@Description: 更新用户信息
//	@receiver u
//	@param ctx
//	@param user
//	@return error
func (u *users) Update(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB {
	db := u.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(user)
}

func newUsers(factory *mysqlFactory) data.UserStore {
	return &users{db: factory.db}
}
