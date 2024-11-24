package glog

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger"
)

const (
	_ginLoggerKey = "_ctx_logger"
	_ginTraceId   = "_ctx_trace_id"
	HeaderTraceId = "X-Trace-traceId"
)

func WithGinCtx(ctx *gin.Context) *gin.Context {
	log := logger.Get().With().Stack().
		Str("traceID", GetTraceId(ctx)).Logger()
	ctx.Set(_ginLoggerKey, &log)
	return ctx
}

func FromGinCtx(ctx *gin.Context) *logger.Logger {
	if l := ctx.Value(_ginLoggerKey); l != nil {
		return l.(*logger.Logger)
	}
	return logger.Get()
}

func Info(ctx *gin.Context, msg string) {
	FromGinCtx(ctx).Info().Msg(msg)
}

func Infof(ctx *gin.Context, format string, v ...interface{}) {
	FromGinCtx(ctx).Info().Msgf(format, v...)
}

func Error(ctx *gin.Context, msg string) {
	FromGinCtx(ctx).Error().Msg(msg)
}

func Errorf(ctx *gin.Context, format string, v ...interface{}) {
	FromGinCtx(ctx).Error().Msgf(format, v...)
}

func Err(ctx *gin.Context, err error, msg string) {
	FromGinCtx(ctx).Err(err).Msg(msg)
}

func Debug(ctx *gin.Context, msg string) {
	FromGinCtx(ctx).Debug().Msg(msg)
}

func Debugf(ctx *gin.Context, format string, v ...interface{}) {
	FromGinCtx(ctx).Debug().Msgf(format, v...)
}

func Warn(ctx *gin.Context, msg string) {
	FromGinCtx(ctx).Warn().Msg(msg)
}

func Warnf(ctx *gin.Context, format string, v ...interface{}) {
	FromGinCtx(ctx).Warn().Msgf(format, v...)
}

func Trace(ctx *gin.Context, msg string) {
	FromGinCtx(ctx).Trace().Msg(msg)
}

func Tracef(ctx *gin.Context, format string, v ...interface{}) {
	FromGinCtx(ctx).Trace().Msgf(format, v...)
}

func GetTraceId(ctx *gin.Context) string {
	if ctx == nil {
		return logger.GenerateLogId()
	}
	if tId := ctx.GetString(_ginTraceId); tId != "" {
		return tId
	}
	// 尝试从header中获取
	var tId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		tId = ctx.GetHeader(HeaderTraceId)
	}
	if tId == "" {
		tId = logger.GenerateLogId()
	}
	ctx.Set(_ginTraceId, tId)
	return tId
}
