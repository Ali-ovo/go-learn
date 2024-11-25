package do

import (
	"shop/pkg/gorm"
)

type CategoryDO struct {
	gorm.BaseModel
	Name             string        `gorm:"type:varchar(20);not null,unique" json:"name"`
	ParentCategoryID int64         `json:"parent"`
	ParentCategory   *CategoryDO   `json:"-"`                                                             // 一对一
	SubCategory      []*CategoryDO `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"` // 一对多
	Level            int32         `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool          `gorm:"default:false;not null" json:"is_tab"`
}

func (CategoryDO) TableName() string {
	return "category"
}

type CategoryDOList struct {
	TotalCount int64         `json:"total,omitempty"`
	Items      []*CategoryDO `json:"data"`
}
