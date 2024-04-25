package render

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/zlog"
	"net/http"
)

type jsonRender struct {
	Errno  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data,omitempty"`
}

func JsonWithSuccess(ctx *gin.Context, data interface{}) {
	r := &jsonRender{
		Errno:  0,
		ErrMsg: "success",
		Data:   data,
	}
	commonHeader(ctx)
	ctx.JSON(http.StatusOK, r)
}

func JsonWithFail(ctx *gin.Context, code int, msg string) {
	r := &jsonRender{
		Errno:  code,
		ErrMsg: msg,
	}
	commonHeader(ctx)
	ctx.JSON(http.StatusOK, r)
}

func JsonWithError(ctx *gin.Context, err error) {
	code, msg := -1, err.Error()
	var err2 Error
	switch {
	case errors.As(err, &err2):
		code = err2.Code
		msg = err2.Msg
	default:
	}
	r := &jsonRender{
		Errno:  code,
		ErrMsg: msg,
	}

	commonHeader(ctx)
	ctx.AbortWithStatusJSON(http.StatusOK, r)
}

func commonHeader(ctx *gin.Context) {
	zlog.SetHeaderLogId(ctx)
}
