package data

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type CategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.CategoryBrandDOList, error)
	GetBrandList(ctx context.Context, categoryID int64) (*do.CategoryBrandDOList, error)
	Get(ctx context.Context, id int64) (*do.CategoryBrandDO, error)
	Create(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) error
	Update(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) error
	Delete(ctx context.Context, txn *gorm.DB, ID uint64) error
}
