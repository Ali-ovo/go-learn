package do

import (
	"shop/pkg/gorm"
)

type BrandsDO struct {
	gorm.BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

func (d *BrandsDO) TableName() string {
	return "brands"
}

type BrandsDOList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*BrandsDO `json:"items,omitempty"`
}
