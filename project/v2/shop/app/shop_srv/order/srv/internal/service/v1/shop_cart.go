package service

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type ShopCartSrv interface {
	Get(ctx context.Context, id int64) (*dto.ShopCartDTO, error)
	List(ctx context.Context, userID int64, meta metav1.ListMeta, orderby []string) (*dto.ShopCartDTOList, error)
	Create(ctx context.Context, shopCart *dto.ShopCartDTO) error
	Update(ctx context.Context, shopCart *dto.ShopCartDTO) error
}
