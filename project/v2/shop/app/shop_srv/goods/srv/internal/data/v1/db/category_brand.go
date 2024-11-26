package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type categoryBrands struct {
	db *gorm.DB
}

func (cb *categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.CategoryBrandDOList, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDOList

	// 这里 赋值是为了保证 db的作用域不受影响
	query := db.Model(&do.CategoryBrandDO{})
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = count

	return &ret, nil
}

func (cb *categoryBrands) GetBrandList(ctx context.Context, categoryID int64) (*do.CategoryBrandDOList, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDOList

	query := db.Model(&do.CategoryBrandDO{})
	query = query.Preload("Category").Preload("Brands").Where("category_id =?", categoryID)
	query.Find(&ret.Items)
	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	ret.TotalCount = query.RowsAffected

	return &ret, nil
}

func (cb *categoryBrands) Get(ctx context.Context, id int64) (*do.CategoryBrandDO, error) {
	db := cb.db.WithContext(ctx)
	var ret do.CategoryBrandDO

	if err := db.Where("id =?", id).First(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryBrandsNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &ret, nil
}

func (cb *categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) error {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Create(gcb)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (cb *categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.CategoryBrandDO) error {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Save(gcb)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (cb *categoryBrands) Delete(ctx context.Context, txn *gorm.DB, ID uint64) error {
	db := cb.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	result := db.Delete(&do.CategoryBrandDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newCategoryBrands(factory *mysqlFactory) *categoryBrands {
	return &categoryBrands{factory.db}
}
