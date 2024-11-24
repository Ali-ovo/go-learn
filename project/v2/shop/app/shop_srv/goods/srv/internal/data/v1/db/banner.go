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

type Banner struct {
	db *gorm.DB
}

func (b *Banner) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerDOList, error) {
	ret := &do.BannerDOList{}

	// 这里 赋值是为了保证 db的作用域不受影响
	query := b.db.Model(&do.BannerDO{})
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = count

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

func newBanner(factory *mysqlFactory) *Banner {
	return &Banner{factory.db}
}

var _ data.BannerStore = (*Banner)(nil)
