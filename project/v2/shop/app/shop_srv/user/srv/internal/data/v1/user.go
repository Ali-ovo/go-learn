package data

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type UserStore interface {
	/*
		有数据访问的方法, 一定要有 error
		参数中最好有 ctx 后期便于管理 比如 cancel掉 链路追踪等
	*/
	// List 用户列表
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.UserDOList, error)
	// GetByMobile 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*do.UserDO, error)
	// GetByID 通过 id 查询用户
	GetByID(ctx context.Context, id uint64) (*do.UserDO, error)
	// Create 添加用户
	Create(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB
	// Update 更新用户
	Update(ctx context.Context, txn *gorm.DB, user *do.UserDO) *gorm.DB
}
