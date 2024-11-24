package srv

import (
	"context"
	"fmt"
	"shop/app/shop_srv/goods/srv/internal/domain/do"
	"shop/app/shop_srv/goods/srv/internal/domain/dto"
	"shop/pkg/gorm"
	"testing"
)

func TestGoods_BatchGet(t *testing.T) {
	dbFactory := conn()

	goodsService := goodsService{
		data: dbFactory,
	}

	get, err := goodsService.BatchGet(context.Background(), []int64{840, 839, 838, 837, 836})
	if err != nil {
		panic(err)
	}

	//get, err = goodsService.BatchGetTwe(context.Background(), []uint64{845, 840, 839, 838, 837, 836})
	//if err != nil {
	//	panic(err)
	//}

	//get, err = goodsService.BatchGetThree(context.Background(), []uint64{845, 840, 839, 838, 837, 836})
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(get)
}

func TestGoods_Update(t *testing.T) {
	dbFactory := conn()

	goodsService := goodsService{
		data: dbFactory,
	}

	goodDO := do.GoodsDO{
		BaseModel:   gorm.BaseModel{ID: 841},
		Name:        "name",
		CategoryID:  135476,
		BrandsID:    686,
		OnSale:      true,
		IsNew:       true,
		IsHot:       true,
		ShipFree:    false,
		GoodsSn:     "asdgfgas",
		ClickNum:    0,
		SoldNum:     0,
		FavNum:      0,
		MarketPrice: 20.5,
		ShopPrice:   33,
		GoodsBrief:  "Lorem eu labore",
		Images: []string{
			"http://hefmcgnjb.dk/krvnlm",
			"http://uvi.ck/clkxg",
			"http://rcsseo.sl/nhybufnrj",
			"http://mbgsqksbx.pl/gky",
			"http://ubr.aq/nhkgjbgdv",
		},
		DescImages: []string{
			"http://dummyimage.com/400x400",
			"http://dummyimage.com/400x400",
			"http://dummyimage.com/400x400",
			"http://dummyimage.com/400x400",
			"http://dummyimage.com/400x400",
		},
		GoodsFrontImage: "http://xbfmxn.kw/uzhs",
	}
	goodDTO := dto.GoodsDTO{goodDO}

	err := goodsService.Update(context.Background(), &goodDTO)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(goodDTO)
}
