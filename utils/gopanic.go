package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger"
)

func WithRecover(ctx *gin.Context, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf(ctx, "goroutine panic recover: %s", err)
		}
	}()

	fn()
}
