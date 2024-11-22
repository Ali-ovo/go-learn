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
	//TODO implement me
	panic("implement me")
}

func (g *goods) Create(ctx context.Context, goods *do.GoodsDO) error {
	//TODO implement me
	panic("implement me")
}

func (g *goods) Update(ctx context.Context, goods *do.GoodsDO) error {
	//TODO implement me
	panic("implement me")
}

func (g *goods) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewGoods(db *gorm.DB) *goods {
	return &goods{db}
}

var _ data.GoodsStore = (*goods)(nil)
