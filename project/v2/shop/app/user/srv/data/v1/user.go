package v1

import (
	"context"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type UserDO struct {
	Name string `json:"name"`
}

type UserDOList struct {
	TotalCount int64     `json:"totalCount,omitempty"` // 总数
	Items      []*UserDO `json:"data"`                 // 数据
}

type UserStore interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDOList, error)
}
