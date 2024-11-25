package dto

import (
	"shop/app/shop/api/internal_srv/domain/do"
)

type CategoryBrandDTO struct {
	do.CategoryBrandDO
}

type CategoryBrandDTOList struct {
	TotalCount int64               `json:"total,omitempty"`
	Items      []*CategoryBrandDTO `json:"data"`
}
