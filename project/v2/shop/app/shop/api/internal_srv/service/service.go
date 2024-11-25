package service

type ServiceFactory interface {
	Banner() BannerSrv
	Brand() BrandSrv
	Category() CategorySrv
	CategoryBrand() CategoryBrandSrv
	Goods() GoodsSrv
}
