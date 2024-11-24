package data

import "gorm.io/gorm"

type DataFactory interface {
	Goods() GoodsStore
	Category() CategoryStore
	Brands() BrandsStore
	Banner() BannerStore
	CategoryBrands() CategoryBrandStore

	Begin() *gorm.DB
}
