package do

import "shop/pkg/gorm"

type CategoryBrandDO struct {
	gorm.BaseModel
	CategoryID int64      `gorm:"type:int;index:idx_category_brand,unique"`
	Category   CategoryDO `gorm:"ForeignKey:CategoryID;reference:ID" json:"category"` // 一对一 设置CategoryBrandDO 表 CategoryID 为外键 与 CategoryDO 的 ID 对应
	BrandsID   int64      `gorm:"type:int;index:idx_category_brand,unique"`
	Brands     BrandsDO   `gorm:"foreignKey:BrandsID;reference:ID" json:"brands"`
}

func (CategoryBrandDO) TableName() string {
	return "goodscategorybrand"
}

type CategoryBrandDOList struct {
	TotalCount int64              `json:"total,omitempty"`
	Items      []*CategoryBrandDO `json:"data,omitempty"`
}
