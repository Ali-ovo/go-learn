package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type categoryBrands struct {
	db *gorm.DB
}

func (cb *categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.CategoryBrandDOList, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	result := db.Model(&do.CategoryBrandDO{})
	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count

	return &ret, nil
}

func (cb *categoryBrands) GetBrandList(ctx context.Context, categoryID int64) (*do.CategoryBrandDOList, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDOList

	result := db.Model(&do.CategoryBrandDO{})
	result = result.Preload("Category").Preload("Brands").Where("category_id =?", categoryID)
	result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = result.RowsAffected

	return &ret, nil
}

func (cb *categoryBrands) Get(ctx context.Context, id int64) (*do.CategoryBrandDO, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDO

	if err := db.Where("id =?", id).First(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (cb *categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) *gorm.DB {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(gcb)
}

func (cb *categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) *gorm.DB {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(gcb)
}

func (cb *categoryBrands) Delete(ctx context.Context, txn *gorm.DB, ID int64) *gorm.DB {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Delete(&do.CategoryBrandDO{}, ID)
}

func newCategoryBrands(factory *mysqlFactory) data.CategoryBrandStore {
	return &categoryBrands{factory.db}
}
