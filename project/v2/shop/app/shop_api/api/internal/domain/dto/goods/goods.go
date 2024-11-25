package dtoGoods

type GoodsFilter struct {
	PriceMin    int32    `form:"pmin" json:"priceMin,omitempty"`
	PriceMax    int32    `form:"pmax" json:"priceMax,omitempty"`
	IsHot       bool     `form:"ih" json:"isHot,omitempty"`
	IsNew       bool     `form:"in" json:"isNew,omitempty"`
	IsTab       bool     `form:"it" json:"isTab,omitempty"`
	TopCategory int32    `form:"c" json:"topCategory,omitempty"`
	Pages       int32    `form:"p" json:"pages,omitempty"`
	PagePerNums int32    `form:"pnum" json:"pagePerNums,omitempty"`
	KeyWords    string   `form:"q" json:"keyWords,omitempty"`
	Brand       int32    `form:"b" json:"brand,omitempty"`
	Orderby     []string `form:"orderby" json:"orderby,omitempty"`
}

type GoodsDTO struct {
	Id              int64        `json:"id,omitempty"`
	CategoryId      int64        `json:"categoryId,omitempty"`
	Name            string       `json:"name,omitempty"`
	GoodsSn         string       `json:"goodsSn,omitempty"`
	ClickNum        int32        `json:"clickNum,omitempty"`
	SoldNum         int32        `json:"soldNum,omitempty"`
	FavNum          int32        `json:"favNum,omitempty"`
	Stocks          int32        `json:"stocks,omitempty"`
	MarketPrice     float32      `json:"marketPrice,omitempty"`
	ShopPrice       float32      `json:"shopPrice,omitempty"`
	GoodsBrief      string       `json:"goodsBrief,omitempty"`
	GoodsDesc       string       `json:"goodsDesc,omitempty"`
	ShipFree        bool         `json:"shipFree,omitempty"`
	Images          []string     `json:"images,omitempty"`
	DescImages      []string     `json:"descImages,omitempty"`
	GoodsFrontImage string       `json:"goodsFrontImage,omitempty"`
	IsNew           bool         `json:"isNew,omitempty"`
	IsHot           bool         `json:"isHot,omitempty"`
	OnSale          bool         `json:"onSale,omitempty"`
	AddTime         int64        `json:"addTime,omitempty"`
	Category        *CategoryDTO `json:"category,omitempty"`
	Brand           *BrandDTO    `json:"brand,omitempty"`
}

type GoodDTOList struct {
	TotalCount int64       `json:"total,omitempty"` // 总数
	Items      []*GoodsDTO `json:"data"`            // 商品数据
}
