package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type category struct {
	db *gorm.DB
}

func (c *category) Get(ctx context.Context, ID int64) (*do.CategoryDO, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDO

	if err := db.Preload("SubCategory.SubCategory").First(&ret, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &ret, nil
}

func (c *category) List(ctx context.Context, level int32) (*do.CategoryDOList, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDOList

	query := db.Where("level =?", level).Find(&ret.Items)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, query.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	query.Count(&ret.TotalCount)
	return &ret, nil
}

func (c *category) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDOList

	query := db.Model(&do.CategoryDO{})
	// 排序
	for _, v := range orderby {
		query = query.Order(v)
	}
	// 加载其他表数据
	query = query.Where("level=1").Preload("SuCategory.SubCategory").Find(&ret.Items)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, query.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	return &ret, query.Error
}

func (c *category) Create(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) error {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Create(category)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *category) Update(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) error {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Updates(category)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *category) Delete(ctx context.Context, txn *gorm.DB, ID int64) error {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	result := db.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newCategory(factory *mysqlFactory) data.CategoryStore {
	return &category{factory.db}
}
