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

type goods struct {
	db *gorm.DB
}

func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	db := g.db.WithContext(ctx)
	var good do.GoodsDO

	if result := db.Preload("Category").Preload("Brands").First(&good, ID); result.RowsAffected == 0 {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return &good, nil
}

func (g *goods) ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	db := g.db.WithContext(ctx)
	var ret do.GoodsDOList

	// 排序
	query := db.Preload("Category").Preload("Brands")
	for _, v := range orderby {
		query = query.Order(v)
	}

	result := query.Where("id in ?", ids).Find(&ret.Items).Count(&ret.TotalCount)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return &ret, nil
}

func (g *goods) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsDOList, error) {
	db := g.db.WithContext(ctx)
	var ret do.GoodsDOList

	// 加载其他表数据
	query := db.Preload("Category").Preload("Brands")
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	result := query.Find(&ret.Items)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	ret.TotalCount = count

	return &ret, nil
}

func (g *goods) Create(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Update(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	tx := db.Updates(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, txn *gorm.DB, ID uint64) error {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}

	result := db.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func newGoods(factory *mysqlFactory) data.GoodsStore {
	return &goods{factory.db}
}
