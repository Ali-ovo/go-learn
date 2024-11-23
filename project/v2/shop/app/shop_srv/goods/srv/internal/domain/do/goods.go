package do

import (
	"database/sql/driver"
	"encoding/json"
	"shop/pkg/gorm"
)

type GoodsSearchDO struct {
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	BrandsID   int32 `json:"brand_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

func (fsd *GoodsSearchDO) GetIndexName() string {
	return "goods"
}

func (fsd *GoodsSearchDO) GetSearchBool() (boolQuery map[string]any, query map[string]any) {
	boolQuery = map[string]any{
		"must":     []any{},
		"must_not": []any{},
		"filter":   []any{},
		"should":   []any{},
	}

	query = map[string]any{
		"query": map[string]any{
			"bool": boolQuery,
		},
		"sort": []any{
			map[string]any{
				"_score": map[string]any{
					"order": "desc",
				},
				"id": map[string]any{
					"order": "asc", // desc 降序
				},
			},
		},
		"from": 0,
		"size": 10,
	}
	return
}

// GetMapping 在 ES 中 POST 一个数据后
// GET goods 获取 生成后的结构 复制过来
func (fsd *GoodsSearchDO) GetMapping() string {
	goodsMapping := `
	{
		"mappings" : {
			"properties" : {
				"brands_id" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"fav_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"ship_free" : {
					"type" : "boolean"
				},
				"shop_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
	return goodsMapping
}

type GoodsSearchDOList struct {
	TotalCount int64            `json:"total_count,omitempty"`
	Items      []*GoodsSearchDO `json:"items"`
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

	Name            string   `gorm:"type:varchar(100);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null"`
	ClickNum        int32    `gorm:"type:int;default:0;not null"`
	SoldNum         int32    `gorm:"type:int;default:0;not null"`
	FavNum          int32    `gorm:"type:int;default:0;not null"`
	MarketPrice     float32  `gorm:"not null"`
	ShopPrice       float32  `gorm:"not null"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`
	Images          GormList `gorm:"type:json;not null"`
	DescImages      GormList `gorm:"type:json;not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}

func (GoodsDO) TableName() string {
	return "goods"
}

type GormList []string

// Value 以json二进制格式 反序列化 保存到数据库中
func (gl *GormList) Value() (driver.Value, error) {
	return json.Marshal(gl)
}

// Scan 以 json 二进制格式 序列化
func (gl *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &gl)
}

type GoodsDOList struct {
	TotalCount int64      `json:"totalCount,omitempty"`
	Items      []*GoodsDO `json:"Items,omitempty"`
}
