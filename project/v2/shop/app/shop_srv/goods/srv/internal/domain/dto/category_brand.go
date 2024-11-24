package dto

import "shop/app/shop_srv/goods/srv/internal/domain/do"

type CategoryBrandDTO struct {
	do.CategoryBrandDO
}

type CategoryBrandDTOList struct {
	TotalCount int64               `json:"total_count,omitempty"`
	Items      []*CategoryBrandDTO `json:"data"`
}
