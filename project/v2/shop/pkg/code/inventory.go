//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrInventoryNotFound - 404: Goods not found.
	ErrInventoryNotFound int = iota + 101201
	// ErrInvSellDetailNotFound - 400: Inventory sell detail not found.
	ErrInvSellDetailNotFound
	// ErrInvNotEnough - 404: Inventory not enough.
	ErrInvNotEnough
	// ErrstockNotFound - 404: Order not found.
	ErrstockNotFound
)
