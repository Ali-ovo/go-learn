//go:generate stringer -type ErrCode -linecomment
package code

type ErrCode int64

const (
	ERR_CODE_OK            ErrCode = 0 // ok
	ERR_CODE_INVALID_PARAM ErrCode = 1 // invalid parameter
	ERR_CODE_TIMEOUT       ErrCode = 2 // timeout
)
