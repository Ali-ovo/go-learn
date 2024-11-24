package dto

import "shop/app/shop_srv/goods/srv/internal/domain/do"

type CategoryDTO struct {
	do.CategoryDO
}

type CategoryDTOList struct {
	TotalCount int64          `json:"total_count,omitempty"`
	Items      []*CategoryDTO `json:"data"`
}
