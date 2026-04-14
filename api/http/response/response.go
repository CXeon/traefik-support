package response

import (
	"errors"
	"net/http"

	tilecontext "github.com/CXeon/tiles/context"
	tileerrors "github.com/CXeon/tiles/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		TraceID: tilecontext.From(c.Request.Context()).TraceID(),
		Data:    data,
	})
}

// Fail 统一处理业务错误响应，HTTP 状态码始终为 200。
// 若 err 是 *tileerrors.Error，直接使用其 Code 和 Message；
// 否则视为未知内部错误，不暴露原始信息。
func Fail(c *gin.Context, err error) {
	traceID := tilecontext.From(c.Request.Context()).TraceID()
	var e *tileerrors.Error
	if errors.As(err, &e) {
		c.JSON(http.StatusOK, Response{
			Code:    e.Code,
			Message: e.Message,
			TraceID: traceID,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    tileerrors.ErrInternal.Code,
		Message: tileerrors.ErrInternal.Message,
		TraceID: traceID,
	})
}
