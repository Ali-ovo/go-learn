package srv

import (
	"shop/app/shop_srv/order/srv/internal/data/v1"
	"shop/app/shop_srv/order/srv/internal/service/v1"
	"shop/pkg/options"
)

type serviceFactory struct {
	data    data.DataFactory
	dtmOpts *options.DtmOptions
}

func (s *serviceFactory) Orders() service.OrderSrv {
	return newOrders(s)
}

func (s *serviceFactory) ShopCart() service.ShopCartSrv {
	panic("eee")
}

func NewService(store data.DataFactory, dtmOptions *options.DtmOptions) service.ServiceFactory {
	return &serviceFactory{data: store, dtmOpts: dtmOptions}
}
