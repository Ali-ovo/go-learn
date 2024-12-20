package do

import "shop/pkg/gorm"

type BrandsDO struct {
	gorm.BaseModel
	Name string `gorm:"type:varchar(50);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

func (d *BrandsDO) TableName() string {
	return "brands"
}

type BrandsDOList struct {
	TotalCount int64       `json:"total,omitempty"`
	Items      []*BrandsDO `json:"data,omitempty"`
}
