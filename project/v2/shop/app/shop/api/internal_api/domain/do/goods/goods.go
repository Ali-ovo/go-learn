package doGoods

type GoodDO struct {
	Id              int64       `json:"id,omitempty"`
	CategoryId      int64       `json:"categoryId,omitempty"`
	Name            string      `json:"name,omitempty"`
	GoodsSn         string      `json:"goodsSn,omitempty"`
	ClickNum        int32       `json:"clickNum,omitempty"`
	SoldNum         int32       `json:"soldNum,omitempty"`
	FavNum          int32       `json:"favNum,omitempty"`
	Stocks          int32       `json:"stocks,omitempty"`
	MarketPrice     float32     `json:"marketPrice,omitempty"`
	ShopPrice       float32     `json:"shopPrice,omitempty"`
	GoodsBrief      string      `json:"goodsBrief,omitempty"`
	GoodsDesc       string      `json:"goodsDesc,omitempty"`
	ShipFree        bool        `json:"shipFree,omitempty"`
	Images          []string    `json:"images,omitempty"`
	DescImages      []string    `json:"descImages,omitempty"`
	GoodsFrontImage string      `json:"goodsFrontImage,omitempty"`
	IsNew           bool        `json:"isNew,omitempty"`
	IsHot           bool        `json:"isHot,omitempty"`
	OnSale          bool        `json:"onSale,omitempty"`
	AddTime         int64       `json:"addTime,omitempty"`
	Category        *CategoryDO `json:"category,omitempty"`
	Brand           *BrandDO    `json:"brand,omitempty"`
}

type GoodDOList struct {
	TotalCount int64     `json:"total_count,omitempty"` // 总数
	Items      []*GoodDO `json:"items"`                 // 商品数据
}
