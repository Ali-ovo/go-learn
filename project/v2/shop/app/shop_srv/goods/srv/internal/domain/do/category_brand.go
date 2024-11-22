package do

import (
	"shop/pkg/gorm"
)

type GoodsCategoryBrandDO struct {
	gorm.BaseModel
	CategoryID int32      `gorm:"type:int;index:idx_category_brand,unique"`
	Category   CategoryDO `gorm:"ForeignKey:CategoryID;reference:ID" json:"category"`

	BrandsID int32    `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   BrandsDO `gorm:"foreignKey:BrandsID;reference:ID" json:"brands"`
}

func (GoodsCategoryBrandDO) TableName() string {
	return "goodscategorybrand"
}

type GoodsCategoryBrandDOList struct {
	TotalCount int64                   `json:"totalCount,omitempty"`
	Items      []*GoodsCategoryBrandDO `json:"Items,omitempty"`
}
