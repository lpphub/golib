package ware

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/zlog"
)

func WithRecover(ctx *gin.Context, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf(ctx, "panic recover: %s", err)
		}
	}()

	fn()
}
