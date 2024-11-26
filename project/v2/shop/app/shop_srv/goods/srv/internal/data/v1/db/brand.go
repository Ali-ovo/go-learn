package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type Brands struct {
	db *gorm.DB
}

func (b *Brands) Get(ctx context.Context, ID int64) (*do.BrandsDO, error) {
	db := b.db.WithContext(ctx)
	var ret do.BrandsDO

	if err := db.Where("id =?", ID).First(ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrBrandsNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &ret, nil
}

func (b *Brands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	db := b.db.WithContext(ctx)
	var ret do.BrandsDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	query := db.Model(&do.BrandsDO{})
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = count

	return &ret, nil
}

func (b *Brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Create(brands)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (b *Brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) (*gorm.DB, error) {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Updates(brands)
	if tx.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return tx, nil
}

func (b *Brands) Delete(ctx context.Context, txn *gorm.DB, ID int64) error {
	db := b.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	result := db.Delete(&do.BrandsDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newBrand(factory *mysqlFactory) data.BrandsStore {
	return &Brands{factory.db}
}
