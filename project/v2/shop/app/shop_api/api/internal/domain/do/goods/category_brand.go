package doGoods

type CategoryBrandDO struct {
	Id         int64 `json:"id,omitempty"`
	CategoryId int64 `json:"categoryId,omitempty"`
	BrandId    int64 `json:"brandId,omitempty"`
}

type CategoryBrandDOList struct {
	TotalCount int64              `json:"total_count,omitempty"` // 总数
	Items      []*CategoryBrandDO `json:"items"`                 // 品牌类型数据
}
