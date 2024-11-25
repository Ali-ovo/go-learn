package db

import (
	"context"
	"shop/app/shop/api/internal_srv/data/v1"
	"shop/app/shop/api/internal_srv/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"

	"gorm.io/gorm"
)

type Banner struct {
	db *gorm.DB
}

func (b *Banner) List(ctx context.Context) (*do.BannerDOList, error) {
	ret := &do.BannerDOList{}

	// 这里 赋值是为了保证 db的作用域不受影响
	query := b.db.Model(&do.BannerDO{})
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = query.RowsAffected

	return ret, nil
}

func (b *Banner) Get(ctx context.Context, id int64) (*do.BannerDO, error) {
	ret := &do.BannerDO{}

	query := b.db.Where("id =?", id).First(&ret)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}

	return ret, nil
}

func (b *Banner) Create(ctx context.Context, banner *do.BannerDO) error {
	tx := b.db.Create(banner)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (b *Banner) Update(ctx context.Context, banner *do.BannerDO) error {
	tx := b.db.Save(banner)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (b *Banner) Delete(ctx context.Context, ID int64) error {
	result := b.db.Delete(&do.BannerDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newBanner(factory *mysqlFactory) data.BannerStore {
	return &Banner{factory.db}
}
