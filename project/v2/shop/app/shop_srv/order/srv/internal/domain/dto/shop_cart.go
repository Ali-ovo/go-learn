package dto

import "shop/app/shop_srv/order/srv/internal/domain/do"

type ShopCartDTO struct {
	do.ShoppingCartDO
}

type ShopCartDTOList struct {
	TotalCount int64          `json:"totalCount,omitempty"`
	Items      []*ShopCartDTO `json:"data"`
}
