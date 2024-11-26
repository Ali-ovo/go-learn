package srv

import (
	"shop/app/shop_srv/user/srv/internal/data/v1"
	"shop/app/shop_srv/user/srv/internal/service"
)

type serviceFactory struct {
	data data.DataFactory
}

func (sf *serviceFactory) User() service.UserSrv {
	return newUser(sf)
}

func NewService(store data.DataFactory) service.ServiceFactory {
	return &serviceFactory{
		data: store,
	}
}
