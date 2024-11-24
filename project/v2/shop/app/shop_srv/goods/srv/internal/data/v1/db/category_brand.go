package db

import (
	"context"

	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"shop/gmicro/pkg/errors"

	"gorm.io/gorm"
)

type categoryBrands struct {
	db *gorm.DB
}

func (cb *categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandDOList, error) {
	ret := &do.GoodsCategoryBrandDOList{}

	// 这里 赋值是为了保证 db的作用域不受影响
	query := cb.db.Model(&do.GoodsCategoryBrandDO{})
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = count

	return ret, nil
}

func (cb *categoryBrands) Create(ctx context.Context, gcb *do.GoodsCategoryBrandDOList) error {
	tx := cb.db.Create(gcb)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (cb *categoryBrands) Update(ctx context.Context, gcb *do.GoodsCategoryBrandDOList) error {
	tx := cb.db.Save(gcb)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (cb *categoryBrands) Delete(ctx context.Context, ID uint64) error {
	result := cb.db.Delete(&do.GoodsCategoryBrandDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newCategoryBrands(factory *mysqlFactory) *categoryBrands {
	return &categoryBrands{factory.db}
}

var _ data.GoodsCategoryBrandStore = (*categoryBrands)(nil)
