package doGoods

type CategoryDO struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	ParentCategory int64  `json:"parentCategory,omitempty"`
	Level          int32  `json:"level,omitempty"`
	IsTab          bool   `json:"isTab,omitempty"`
}

type CategoryDOList struct {
	TotalCount int64         `json:"total_count,omitempty"` // 总数
	Items      []*CategoryDO `json:"items"`                 // 类型数据
}
