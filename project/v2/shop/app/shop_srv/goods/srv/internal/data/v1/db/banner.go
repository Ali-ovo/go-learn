package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"

	"gorm.io/gorm"
)

type Banner struct {
	db *gorm.DB
}

func (b *Banner) List(ctx context.Context) (*do.BannerDOList, error) {
	db := b.db.WithContext(ctx)
	var ret do.BannerDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.BannerDO{})
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = result.RowsAffected

	return &ret, nil
}

func (b *Banner) Get(ctx context.Context, id int64) (*do.BannerDO, error) {
	db := b.db.WithContext(ctx)
	var ret do.BannerDO

	result := db.Where("id =?", id).First(&ret)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ret, nil
}

func (b *Banner) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(banner)
}

func (b *Banner) Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(banner)
}

func (b *Banner) Delete(ctx context.Context, txn *gorm.DB, ID int64) *gorm.DB {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Delete(&do.BannerDO{}, ID)
}

func newBanner(factory *mysqlFactory) data.BannerStore {
	return &Banner{factory.db}
}
