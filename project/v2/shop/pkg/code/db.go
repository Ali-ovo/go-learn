//go:generate codegen -type=int -doc -output ./error_code_generated.md
//go:generate codegen -type=int
package code

const (
	// ErrConnectDB - 500: init db error.
	ErrConnectDB int = iota + 100501
)

const (
	// ErrEsDatabase - 404: EsDatabase error.
	ErrEsDatabase int = iota + 100601
	// ErrEsUnmarshal - 500: Es unmarshal error.
	ErrEsUnmarshal
)

const (
	// ErrConnectGRPC - 500: Connect to grpc error.
	ErrConnectGRPC int = iota + 100701
)

const (
	// ErrRedisDatabase - 500: Redis data base error.
	ErrRedisDatabase int = iota + 100801
)
