package service

import (
	"context"
	"shop/app/shop_srv/user/srv/internal/domain/dto"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type UserSrv interface {
	// List 获取 用户列表页
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*dto.UserDTOList, error)
	// GetByMobile 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*dto.UserDTO, error)
	// GetByID 通过 id 查询用户
	GetByID(ctx context.Context, id uint64) (*dto.UserDTO, error)
	// Create 添加用户
	Create(ctx context.Context, user *dto.UserDTO) error
	// Update 更新用户
	Update(ctx context.Context, user *dto.UserDTO) error
}
