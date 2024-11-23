//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 101101

	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound

	// ErrBrandsNotFound - 404: Brand not found.
	ErrBrandsNotFound
)
