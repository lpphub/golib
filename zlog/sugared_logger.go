package zlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	sugaredLoggerAddr = "_sugared_addr"
)

func GetLogger(lc LogConf) *zap.SugaredLogger {
	if SugaredLogger == nil {
		if ZapLogger == nil {
			ZapLogger = newLogger(lc).WithOptions(zap.AddCallerSkip(1))
		}
		SugaredLogger = ZapLogger.Sugar()
	}
	return SugaredLogger
}

func sugaredLogger(ctx *gin.Context) *zap.SugaredLogger {
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
	sugaredLogger(ctx).Debug(args...)
}

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Debugf(format, args...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Info(args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Infof(format, args...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Warn(args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Error(args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Errorf(format, args...)
}

func Fatal(ctx *gin.Context, args ...interface{}) {
	sugaredLogger(ctx).Fatal(args...)
}

func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
	sugaredLogger(ctx).Fatalf(format, args...)
}
