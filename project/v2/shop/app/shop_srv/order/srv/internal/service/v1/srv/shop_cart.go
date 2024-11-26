package srv

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	"shop/app/shop_srv/order/srv/internal/service/v1"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type ShopCartSrv struct {
	data data.DataFactory
}

func (s ShopCartSrv) Get(ctx context.Context, id int64) (*dto.ShopCartDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (s ShopCartSrv) List(ctx context.Context, userID int64, meta metav1.ListMeta, orderby []string) (*dto.ShopCartDTOList, error) {
	//TODO implement me
	panic("implement me")
}

func (s ShopCartSrv) Create(ctx context.Context, shopCart *dto.ShopCartDTO) error {
	//TODO implement me
	panic("implement me")
}

func (s ShopCartSrv) Update(ctx context.Context, shopCart *dto.ShopCartDTO) error {
	//TODO implement me
	panic("implement me")
}

func newShopCart(srv *serviceFactory) service.ShopCartSrv {
	return &ShopCartSrv{
		data: srv.data,
	}
}
