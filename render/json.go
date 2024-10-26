package render

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/zlog"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type JsonRender struct {
	Errno  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data,omitempty"`
}

func JsonWithSuccess(ctx *gin.Context, data interface{}) {
	r := &JsonRender{
		Errno:  0,
		ErrMsg: "success",
		Data:   data,
	}
	commonHeader(ctx)
	ctx.JSON(http.StatusOK, r)
}

func JsonWithError(ctx *gin.Context, err error) {
	code, msg := -1, err.Error()

	var err2 Error
	if errors.As(err, &err2) {
		code = err2.Code
		msg = err2.Msg
	}

	r := &JsonRender{
		Errno:  code,
		ErrMsg: msg,
	}
	commonHeader(ctx)
	ctx.AbortWithStatusJSON(http.StatusOK, r)
}

func JsonWithFail(ctx *gin.Context, code int, msg string) {
	r := &JsonRender{
		Errno:  code,
		ErrMsg: msg,
	}
	commonHeader(ctx)
	ctx.JSON(http.StatusOK, r)
}

func JsonAbortWithFail(ctx *gin.Context, code int, msg string) {
	ctx.Abort()
	JsonWithFail(ctx, code, msg)
}

func commonHeader(ctx *gin.Context) {
	zlog.SetHeaderLogId(ctx)
	ctx.Header("X-Resp-Time", fmt.Sprintf("%d", time.Now().Unix()))
}
