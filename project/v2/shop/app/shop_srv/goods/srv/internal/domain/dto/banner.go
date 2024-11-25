package dto

import "shop/app/shop_srv/goods/srv/internal/domain/do"

type BannerDTO struct {
	do.BannerDO
}

type BannerDTOList struct {
	TotalCount int64        `json:"total,omitempty"`
	Items      []*BannerDTO `json:"data"`
}
