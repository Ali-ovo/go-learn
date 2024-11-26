package do

import "shop/pkg/gorm"

type BannerDO struct {
	gorm.BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

func (BannerDO) TableName() string {
	return "banner"
}

type BannerDOList struct {
	TotalCount int64       `json:"total,omitempty"`
	Items      []*BannerDO `json:"data"`
}
