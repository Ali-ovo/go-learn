package doGoods

type BannerDO struct {
	Id    int64  `json:"id,omitempty"`
	Index int32  `json:"index,omitempty"`
	Image string `json:"image,omitempty"`
	Url   string `json:"url,omitempty"`
}

type BannerDOList struct {
	TotalCount int64       `json:"total_count,omitempty"` // 总数
	Items      []*BannerDO `json:"items"`                 // 轮播图数据
}
