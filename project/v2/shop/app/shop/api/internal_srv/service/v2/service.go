package srv

import (
	"shop/app/shop/api/internal_srv/data/v1"
	"shop/app/shop/api/internal_srv/data_search/v1"
	"shop/app/shop/api/internal_srv/service"
)

type serviceFactory struct {
	data       data.DataFactory
	dataSearch data_search.SearchFactory
}

func (sf *serviceFactory) Banner() service.BannerSrv {
	return newBanner(sf)
}

func (sf *serviceFactory) Brand() service.BrandSrv {
	return newBrand(sf)
}

func (sf *serviceFactory) Category() service.CategorySrv {
	return newCategory(sf)
}

func (sf *serviceFactory) CategoryBrand() service.CategoryBrandSrv {
	return newCategoryBrand(sf)
}

func (sf *serviceFactory) Goods() service.GoodsSrv {
	return newGoods(sf)
}

func NewService(store data.DataFactory, dataSearch data_search.SearchFactory) service.ServiceFactory {
	return &serviceFactory{data: store, dataSearch: dataSearch}
}
