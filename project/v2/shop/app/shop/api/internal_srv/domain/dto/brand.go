package dto

import (
	"shop/app/shop/api/internal_srv/domain/do"
)

type BrandsDTO struct {
	do.BrandsDO
}

type BrandsDTOList struct {
	TotalCount int64        `json:"total,omitempty"`
	Items      []*BrandsDTO `json:"data"`
}
