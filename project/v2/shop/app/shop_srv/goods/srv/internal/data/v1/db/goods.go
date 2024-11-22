package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/gmicro/pkg/code"
	"shop/gmicro/pkg/errors"

	metav1 "shop/gmicro/pkg/common/meta/v1"
	code2 "shop/pkg/code"

	"gorm.io/gorm"
)

type goods struct {
	db *gorm.DB
}

func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	good := do.GoodsDO{}
	if result := g.db.Preload("Category").Preload("Brands").First(&good, ID); result.RowsAffected == 0 {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return &good, nil
}

func (g *goods) ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}

	// 排序
	query := g.db.Preload("Category").Preload("Brands")
	for _, v := range orderby {
		query = query.Order(v)
	}

	d := query.Where("id in ?", ids).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsDOList, error) {
	ret := &do.GoodsDOList{}

	// 处理分页
	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}
	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}

	// 排序
	query := g.db.Preload("Category").Preload("Brands")
	for _, v := range orderby {
		query = query.Order(v)
	}

	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, d.Error.Error())
	}
	return ret, nil
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}

	return nil
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	result := g.db.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func NewGoods(db *gorm.DB) *goods {
	return &goods{db}
}

var _ data.GoodsStore = (*goods)(nil)
