package do

import (
	"database/sql/driver"
	"encoding/json"
	"shop/pkg/gorm"
)

type GormList []string

func (gl GormList) Value() (driver.Value, error) {
	return json.Marshal(gl)
}

type GoodsDO struct {
	gorm.BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   CategoryDO
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     BrandsDO

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}

func (GoodsDO) TableName() string {
	return "goods"
}

type GoodsDOList struct {
	TotalCount int64      `json:"totalCount,omitempty"`
	Items      []*GoodsDO `json:"Items,omitempty"`
}
