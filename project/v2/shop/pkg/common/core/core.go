package core

import (
	"fmt"
	"go-learn/project/v2/shop/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"msg"`

	Detail string `json:"detail"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
// WriteResponse将错误或响应数据写入HTTP响应体中。
// 它使用errors.ParseCoder将任何错误解析为errors.Coder。
// errors.Coder包含错误代码、用户友好的错误消息和HTTP状态码。
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		errStr := fmt.Sprintf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Detail:    errStr,
			Reference: coder.Reference(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
