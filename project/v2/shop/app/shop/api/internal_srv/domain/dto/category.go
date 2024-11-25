package dto

import (
	"shop/app/shop/api/internal_srv/domain/do"
)

type CategoryDTO struct {
	do.CategoryDO
}

type CategoryDTOList struct {
	TotalCount int64          `json:"total,omitempty"`
	Items      []*CategoryDTO `json:"data"`
}
