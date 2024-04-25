package zlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	sugaredLoggerAddr = "_sugared_addr"
)

func sugaredLoggerWithCtx(ctx *gin.Context) *zap.SugaredLogger {
	if ctx == nil {
		return SugaredLogger
	}
	if t, exist := ctx.Get(sugaredLoggerAddr); exist {
		if s, ok := t.(*zap.SugaredLogger); ok {
			return s
		}
	}
	s := SugaredLogger.With(
		zap.String("logId", GetLogId(ctx)),
		zap.String("url", ctx.Request.URL.Path),
	)
	ctx.Set(sugaredLoggerAddr, s)
	return s
}

func Debug(ctx *gin.Context, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Debug(args...)
}

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Debugf(format, args...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Info(args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Infof(format, args...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Warn(args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Error(args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Errorf(format, args...)
}

func Fatal(ctx *gin.Context, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Fatal(args...)
}

func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLoggerWithCtx(ctx).Fatalf(format, args...)
}
