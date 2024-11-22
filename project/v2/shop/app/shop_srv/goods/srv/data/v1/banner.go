package v1

import (
	"context"
	"shop/pkg/gorm"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BannerDO struct {
	gorm.BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

func (BannerDO) TableName() string {
	return "banner"
}

type BannerDOList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*BannerDO `json:"data"`
}

type BannerData interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*BannerDOList, error)
	Create(ctx context.Context, banner *BannerDO) error
	Update(ctx context.Context, banner *BannerDO) error
	Delete(ctx context.Context, ID int64) error
}
