package dto

import "shop/app/shop_srv/goods/srv/internal/domain/do"

type BrandsDTO struct {
	do.BrandsDO
}

type BrandsDTOList struct {
	TotalCount int64        `json:"total,omitempty"`
	Items      []*BrandsDTO `json:"data"`
}
