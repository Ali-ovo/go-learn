package dto

import "shop/app/shop_srv/user/srv/internal/domain/do"

type UserDTO struct {
	// 这里偷个懒, 应为业务层和 底层 字段没有太大变动
	do.UserDO
}

type UserDTOList struct {
	TotalCount int64      `json:"totalCount,omitempty"` // 总数
	Items      []*UserDTO `json:"data"`                 // 数据
}
