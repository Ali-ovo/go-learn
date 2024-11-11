//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100501
)
