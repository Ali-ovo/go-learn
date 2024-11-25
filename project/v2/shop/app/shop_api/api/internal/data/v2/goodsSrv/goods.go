package goodsSrv

import (
	Igoods "shop/app/shop_api/api/internal/data/v2/goodsSrv/goods"
)

type GoodsData interface {
	Goods() Igoods.Goods
	Brand() Igoods.Brand
}
