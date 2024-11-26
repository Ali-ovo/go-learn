//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrOrderNotFound - 404: Order not found.
	ErrOrderNotFound int = iota + 101301
	// ErrShopCartNotFound - 404: ShopCart not found.
	ErrShopCartNotFound
	// ErrOrderDtm - 404: Dtm unknonwn error.
	ErrOrderDtm
	// ErrNotGoodsSelect - 404: No Goods selected.
	ErrNotGoodsSelect
)
