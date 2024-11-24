package srv

import (
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/data_search/v1"
	"shop/app/shop_srv/goods/srv/internal/service"
)

type serviceFactory struct {
	data       data.DataFactory
	dataSearch data_search.SearchFactory
}

func (sf *serviceFactory) Goods() service.GoodsSrv {
	return newGoods(sf)
}

func NewService(store data.DataFactory, dataSearch data_search.SearchFactory) *serviceFactory {
	return &serviceFactory{data: store, dataSearch: dataSearch}
}

var _ service.ServiceFactory = (*serviceFactory)(nil)
