package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger/logx"
)

func WithRecover(ctx *gin.Context, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logx.Errorf(ctx, "goroutine panic recover: %s", err)
		}
	}()

	fn()
}
