package dto

import "shop/app/shop_srv/goods/srv/internal/domain/do"

type GoodsDTO struct {
	do.GoodsDO
}

type GoodsDTOList struct {
	TotalCount int64       `json:"total_count,omitempty"`
	Items      []*GoodsDTO `json:"data"`
}
