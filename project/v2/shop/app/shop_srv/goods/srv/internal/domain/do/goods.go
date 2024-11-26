package do

import (
	"database/sql/driver"
	"encoding/json"
	"shop/pkg/gorm"
	"time"
)

type GoodsSearchDO struct {
	ID          int64      `json:"id" mapstructure:"id"`
	CategoryID  int64      `json:"category_id" mapstructure:"category_id"`
	BrandsID    int64      `json:"brand_id" mapstructure:"brand_id"`
	OnSale      bool       `json:"on_sale" mapstructure:"on_sale"`
	ShipFree    bool       `json:"ship_free" mapstructure:"ship_free"`
	IsNew       bool       `json:"is_new" mapstructure:"is_new"`
	IsHot       bool       `json:"is_hot" mapstructure:"is_hot"`
	Name        string     `json:"name" mapstructure:"name"`
	ClickNum    int32      `json:"click_num" mapstructure:"click_num"`
	SoldNum     int32      `json:"sold_num" mapstructure:"sold_num"`
	FavNum      int32      `json:"fav_num" mapstructure:"fav_num"`
	MarketPrice float32    `json:"market_price" mapstructure:"market_price"`
	GoodsBrief  string     `json:"goods_brief" mapstructure:"goods_brief"`
	ShopPrice   float32    `json:"shop_price" mapstructure:"shop_price"`
	DeleteAt    *time.Time `json:"delete_at" mapstructure:"delete_at"`
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
		"sort": []map[string]any{
			{
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

	CategoryID int64 `gorm:"type:int;not null" json:"category_id"`
	Category   CategoryDO
	BrandsID   int64 `gorm:"type:int;not null" json:"brands_id"`
	Brands     BrandsDO

	OnSale   bool `gorm:"default:false;not null" json:"on_sale"`
	ShipFree bool `gorm:"default:false;not null" json:"ship_free"`
	IsNew    bool `gorm:"default:false;not null" json:"is_new"`
	IsHot    bool `gorm:"default:false;not null" json:"is_hot"`

	Name            string   `gorm:"type:varchar(100);not null" json:"name"`
	GoodsSn         string   `gorm:"type:varchar(50);not null" json:"goods_sn"`    // 商品编号
	ClickNum        int32    `gorm:"type:int;default:0;not null" json:"click_num"` // 点击次数	TODO 单独创建一张表
	SoldNum         int32    `gorm:"type:int;default:0;not null" json:"sold_num"`  // 购买次数
	FavNum          int32    `gorm:"type:int;default:0;not null" json:"fav_num"`   // 收藏次数
	MarketPrice     float32  `gorm:"not null" json:"market_price"`
	ShopPrice       float32  `gorm:"not null" json:"shop_price"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null" json:"goods_brief"`
	Images          GormList `gorm:"type:json;not null" json:"images"`
	DescImages      GormList `gorm:"type:json;not null" json:"desc_images"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null" json:"goods_front_image"`
}

func (GoodsDO) TableName() string {
	return "goods"
}

type GormList []string

// Value 以json二进制格式 反序列化 保存到数据库中
func (gl GormList) Value() (driver.Value, error) {
	tmp, err := json.Marshal(gl)
	return tmp, err
}

// Scan 以 json 二进制格式 序列化
func (gl *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &gl)
}

type GoodsDOList struct {
	TotalCount int64      `json:"total,omitempty"`
	Items      []*GoodsDO `json:"data,omitempty"`
}
