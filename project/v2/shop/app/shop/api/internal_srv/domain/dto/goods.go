package dto

import (
	"shop/app/shop/api/internal_srv/domain/do"
)

type GoodsDTO struct {
	do.GoodsDO
}

type GoodsDTOList struct {
	TotalCount int64       `json:"total,omitempty"`
	Items      []*GoodsDTO `json:"data"`
}
