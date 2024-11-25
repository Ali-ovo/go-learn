package doGoods

type BrandDO struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Logo string `json:"logo,omitempty"`
}

type BrandDOList struct {
	TotalCount int64      `json:"total_count,omitempty"` // 总数
	Items      []*BrandDO `json:"items"`                 // 品牌数据
}
