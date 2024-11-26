package db

import (
	"context"
	"shop/app/shop_srv/goods/srv/internal/data/v1"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	metav1 "shop/gmicro/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type goods struct {
	db *gorm.DB
}

func (g *goods) Get(ctx context.Context, ID uint64) (*do.GoodsDO, error) {
	db := g.db.WithContext(ctx)
	var good do.GoodsDO

	if result := db.Preload("Category").Preload("Brands").First(&good, ID); result.RowsAffected == 0 {
		return nil, result.Error
	}
	return &good, nil
}

func (g *goods) ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error) {
	db := g.db.WithContext(ctx)
	var ret do.GoodsDOList

	// 排序
	result := db.Preload("Category").Preload("Brands")
	for _, v := range orderby {
		result = result.Order(v)
	}

	result = result.Where("id in ?", ids).Find(&ret.Items).Count(&ret.TotalCount)
	if result.Error != nil {
		return nil, result.Error
	}
	return &ret, nil
}

func (g *goods) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsDOList, error) {
	db := g.db.WithContext(ctx)
	var ret do.GoodsDOList

	// 加载其他表数据
	result := db.Preload("Category").Preload("Brands")
	// 处理分页 排序
	result, count := paginate(result, opts.Page, opts.PageSize, orderby)
	result = result.Find(&ret.Items)
	if result.Error != nil {
		return nil, result.Error
	}
	ret.TotalCount = count

	return &ret, nil
}

func (g *goods) Create(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) *gorm.DB {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Create(goods)
}

func (g *goods) Update(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) *gorm.DB {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Updates(goods)
}

func (g *goods) Delete(ctx context.Context, txn *gorm.DB, ID uint64) *gorm.DB {
	db := g.db.WithContext(ctx)
	if txn != nil {
		db = txn.WithContext(ctx)
	}
	return db.Delete(&do.GoodsDO{}, ID)
}

func newGoods(factory *mysqlFactory) data.GoodsStore {
	return &goods{factory.db}
}
