package dtoGoods

type CategoryBrandDTO struct {
	Id         int64 `json:"id,omitempty"`
	CategoryId int64 `json:"categoryId,omitempty"`
	BrandId    int64 `json:"brandId,omitempty"`
}

type CategoryBrandDTOList struct {
	TotalCount int64               `json:"total_count,omitempty"` // 总数
	Items      []*CategoryBrandDTO `json:"items"`                 // 品牌类型数据
}
