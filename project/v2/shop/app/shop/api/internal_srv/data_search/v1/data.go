package data_search

type SearchFactory interface {
	Goods() GoodsStore
}
