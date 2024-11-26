package do

import "shop/pkg/gorm"

type ShoppingCartDO struct {
	gorm.BaseModel
	User    int32 `gorm:"type:int;index"` // 在购物车列表中我们需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index"` // 加索引: 需要查询的时候 添加 1. 会影响插入性能 2. 会占用磁盘
	Nums    int32 `gorm:"type:int"`       // 购买的数量
	Checked bool  // 是否选中
}

func (ShoppingCartDO) TableName() string {
	return "shoppingcart"
}

type ShoppingCartDOList struct {
	TotalCount int64             `json:"totalCount,omitempty"`
	Items      []*ShoppingCartDO `json:"items"`
}
