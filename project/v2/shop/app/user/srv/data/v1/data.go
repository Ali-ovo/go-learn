package v1

import (
	"context"
	metav1 "go-learn/project/v2/shop/pkg/common/meta/v1"
)

type UserDO struct {
	Name string `json:"name"`
}

type UserDOList struct {
	TotalCount int64     `json:"totalCount,omitempty"`
	Items      []*UserDO `json:"data"`
}

type UserStore interface {
	List(ctx context.Context, opts metav1.ListMeta) (*UserDOList, error)
}
