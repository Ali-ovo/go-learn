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

func (c *category) Get(ctx context.Context, ID int32) (*do.CategoryDO, error) {
	ret := &do.CategoryDO{}

	if err := c.db.Preload("SubCategory.SubCategory").First(&ret, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrCategoryNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return ret, nil
}

func (c *category) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	ret := &do.CategoryDOList{}
	query := c.db.Model(&do.CategoryDO{})

	// 排序
	for _, v := range orderby {
		query = query.Order(v)
	}
	// 加载其他表数据
	query = query.Where("level=1").Preload("SuCategory.SubCategory").Find(&ret.Items)

	return ret, query.Error
}

func (c *category) Create(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Create(category)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *category) Update(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Save(category)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *category) Delete(ctx context.Context, ID int64) error {
	result := c.db.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func NewCategory(db *gorm.DB) *category {
	return &category{db}
}

var _ data.CategoryStore = (*category)(nil)
