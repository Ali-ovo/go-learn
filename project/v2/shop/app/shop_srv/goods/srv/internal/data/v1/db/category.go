package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"

	"gorm.io/gorm"
)

type category struct {
	db *gorm.DB
}

func (c *category) Get(ctx context.Context, ID int64) (*do.CategoryDO, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDO

	if err := db.Preload("SubCategory.SubCategory").First(&ret, ID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *category) List(ctx context.Context, level int32) (*do.CategoryDOList, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDOList

	result := db.Where("level =?", level).Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	result.Count(&ret.TotalCount)
	return &ret, nil
}

func (c *category) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	db := c.db.WithContext(ctx)
	var ret do.CategoryDOList

	result := db.Model(&do.CategoryDO{})
	// 排序
	for _, v := range orderby {
		result = result.Order(v)
	}
	// 加载其他表数据
	result = result.Where("level=1").Preload("SuCategory.SubCategory").Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	result.Count(&ret.TotalCount)
	return &ret, nil
}

func (c *category) Create(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) *gorm.DB {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(category)
}

func (c *category) Update(ctx context.Context, txn *gorm.DB, category *do.CategoryDO) *gorm.DB {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(category)
}

func (c *category) Delete(ctx context.Context, txn *gorm.DB, ID int64) *gorm.DB {
	db := c.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Delete(&do.GoodsDO{}, ID)
}

func newCategory(factory *mysqlFactory) data.CategoryStore {
	return &category{factory.db}
}
