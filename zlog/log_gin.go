package zlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	_ginLogger = "_ctx_logger"
)

func loggerWithGinCtx(ctx *gin.Context) *zap.SugaredLogger {
	if ctx == nil {
		return SugaredLogger
	}
	if t, exist := ctx.Get(_ginLogger); exist {
		if s, ok := t.(*zap.SugaredLogger); ok {
			return s
		}
	}

	s := SugaredLogger.With(
		zap.String("logId", GetLogId(ctx)),
		zap.String("module", GetModuleWithDefault(ctx, "app")),
	)
	ctx.Set(_ginLogger, s)
	return s
}

func Debug(ctx *gin.Context, args ...interface{}) {
	loggerWithGinCtx(ctx).Debug(args...)
}

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	loggerWithGinCtx(ctx).Debugf(format, args...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	loggerWithGinCtx(ctx).Info(args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	loggerWithGinCtx(ctx).Infof(format, args...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	loggerWithGinCtx(ctx).Warn(args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	loggerWithGinCtx(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	loggerWithGinCtx(ctx).Error(args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	loggerWithGinCtx(ctx).Errorf(format, args...)
}

func Fatal(ctx *gin.Context, args ...interface{}) {
	loggerWithGinCtx(ctx).Fatal(args...)
}

func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
	loggerWithGinCtx(ctx).Fatalf(format, args...)
}
