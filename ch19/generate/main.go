//go:generate stringer -type ErrCode -linecomment
package main

import (
	"fmt"
	"go-learn/ch19/generate/code"
)

type ErrCode int64

const (
	ERR_CODE_OK            ErrCode = 0 // ok
	ERR_CODE_INVALID_PARAM ErrCode = 1 // invalid parameter
	ERR_CODE_TIMEOUT       ErrCode = 2 // timeout
)

// var mapErrDesc = map[int]string{
// 	ERR_CODE_OK:            "ok",
// 	ERR_CODE_INVALID_PARAM: "invalid parameter",
// 	ERR_CODE_TIMEOUT:       "timeout",
// }

// func GetErrDesc(code int) string {
// 	if desc, ok := mapErrDesc[code]; ok {
// 		return desc
// 	}
// 	return "unknown error"
// }

func main() {
	fmt.Println(code.ERR_CODE_INVALID_PARAM)
}
