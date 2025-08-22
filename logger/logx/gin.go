package logx

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger"
)

const (
	_ginLogger  = "_ctx_logger"
	_ginLogId   = "_ctx_log_id"
	HeaderLogId = "X-Trace-logId"
)

func FromGinCtx(ctx *gin.Context) *logger.Logger {
	if ctx == nil {
		return logger.Log()
	}
	if l := ctx.Value(_ginLogger); l != nil {
		return l.(*logger.Logger)
	}

	log := logger.Log().With().CallerWithSkipFrameCount(3).Str("logId", GetLogId(ctx)).Logger()
	ctx.Set(_ginLogger, &log)
	return &log
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

func GetLogId(ctx *gin.Context) string {
	if ctx == nil {
		return logger.GenerateLogId()
	}
	if tId := ctx.GetString(_ginLogId); tId != "" {
		return tId
	}
	// 尝试从header中获取
	var tId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		tId = ctx.GetHeader(HeaderLogId)
	}
	if tId == "" {
		tId = logger.GenerateLogId()
	}
	ctx.Set(_ginLogId, tId)
	return tId
}
