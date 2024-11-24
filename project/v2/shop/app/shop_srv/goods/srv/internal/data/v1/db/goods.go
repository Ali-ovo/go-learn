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

	// 加载其他表数据
	query := g.db.Preload("Category").Preload("Brands")
	// 处理分页 排序
	query, count := paginate(query, opts.Page, opts.PageSize, orderby)
	query.Find(&ret.Items)
	if query.Error != nil {
		return nil, errors.WithCode(code.ErrDatabase, query.Error.Error())
	}
	ret.TotalCount = count

	return ret, nil
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsDO) error {
	tx := g.db.Create(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Create(goods)
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

func (g *goods) UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error {
	tx := txn.Save(goods)
	if tx.Error != nil {
		return errors.WithCode(code.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	result := g.db.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func (g *goods) DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error {
	result := txn.Delete(&do.GoodsDO{}, ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.WithCode(code2.ErrGoodsNotFound, result.Error.Error())
		}
		return errors.WithCode(code.ErrDatabase, result.Error.Error())
	}
	return nil
}

func (g *goods) Begin(ctx context.Context) *gorm.DB {
	return g.db.Begin()
}

func NewGoods(db *gorm.DB) *goods {
	return &goods{db}
}

var _ data.GoodsStore = (*goods)(nil)
