package dtoGoods

type CategoryDTO struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	ParentCategory int64  `json:"parentCategory,omitempty"`
	Level          int32  `json:"level,omitempty"`
	IsTab          bool   `json:"isTab,omitempty"`
}

type CategoryDTOList struct {
	TotalCount int64          `json:"total_count,omitempty"` // 总数
	Items      []*CategoryDTO `json:"items"`                 // 类型数据
}
