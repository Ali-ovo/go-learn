package v1

import (
	"context"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/pkg/gorm"
)

type BrandsDO struct {
	gorm.BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

func (d *BrandsDO) TableName() string {
	return "brands"
}

type BrandsDOList struct {
	TotalCount int64       `json:"totalCount,omitempty"`
	Items      []*BrandsDO `json:"items,omitempty"`
}

type BrandsStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*BrandsDOList, error)
	Create(ctx context.Context, brands *BrandsDO) error
	Update(ctx context.Context, brands *BrandsDO) error
	Delete(ctx context.Context, ID uint64) error
}
