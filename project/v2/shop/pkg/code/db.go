//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrConnectDB - 500: init db error.
	ErrConnectDB int = iota + 100601
)
