package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"

	"gorm.io/gorm"
)

type Banner struct {
	db *gorm.DB
}

func (b *Banner) List(ctx context.Context) (*do.BannerDOList, error) {
	db := b.db.WithContext(ctx)
	var ret do.BannerDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	query := db.Model(&do.BannerDO{})
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = query.RowsAffected

	return &ret, nil
}

func (b *Banner) Get(ctx context.Context, id int64) (*do.BannerDO, error) {
	db := b.db.WithContext(ctx)
	var ret do.BannerDO

	query := db.Where("id =?", id).First(&ret)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}

	return &ret, nil
}

func (b *Banner) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Create(banner)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (b *Banner) Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Save(banner)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (b *Banner) Delete(ctx context.Context, txn *gorm.DB, ID int64) error {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	result := db.Delete(&do.BannerDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newBanner(factory *mysqlFactory) data.BannerStore {
	return &Banner{factory.db}
}
