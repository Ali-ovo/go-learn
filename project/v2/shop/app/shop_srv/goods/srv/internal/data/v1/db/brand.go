package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type Brands struct {
	db *gorm.DB
}

func (b *Brands) Get(ctx context.Context, ID int64) (*do.BrandsDO, error) {
	db := b.db.WithContext(ctx)
	var ret do.BrandsDO

	if err := db.Where("id =?", ID).First(ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *Brands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	db := b.db.WithContext(ctx)
	var ret do.BrandsDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.BrandsDO{})
	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count

	return &ret, nil
}

func (b *Brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(brands)
}

func (b *Brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(brands)
}

func (b *Brands) Delete(ctx context.Context, txn *gorm.DB, ID int64) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Delete(&do.BrandsDO{}, ID)
}

func newBrand(factory *mysqlFactory) data.BrandsStore {
	return &Brands{factory.db}
}
