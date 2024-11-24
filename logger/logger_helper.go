package logger

import (
	"context"
	"github.com/rs/zerolog"
	"strconv"
	"time"
)

const (
	CtxTraceID = "ctx_traceId"
)

func WithCtx(ctx context.Context) context.Context {
	traceID := ""
	if tid := ctx.Value(CtxTraceID); tid != nil {
		traceID = tid.(string)
	} else {
		traceID = GenerateLogId()
	}

	log := logger.With().Stack().Str("traceID", traceID).Logger()
	return log.WithContext(ctx)
}

func FromCtx(ctx context.Context) *Logger {
	if l := zerolog.Ctx(ctx); l != nil && l.GetLevel() != zerolog.Disabled {
		return l
	}
	return logger
}

func Info(ctx context.Context, msg string) {
	FromCtx(ctx).Info().Msg(msg)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Info().Msgf(format, v...)
}

func Error(ctx context.Context, msg string) {
	FromCtx(ctx).Error().Msg(msg)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Error().Msgf(format, v...)
}

func Err(ctx context.Context, err error, msg string) {
	FromCtx(ctx).Err(err).Msg(msg)
}

func Debug(ctx context.Context, msg string) {
	FromCtx(ctx).Debug().Msg(msg)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Debug().Msgf(format, v...)
}

func Warn(ctx context.Context, msg string) {
	FromCtx(ctx).Warn().Msg(msg)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Warn().Msgf(format, v...)
}

func Trace(ctx context.Context, msg string) {
	FromCtx(ctx).Trace().Msg(msg)
}

func Tracef(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Trace().Msgf(format, v...)
}

func GenerateLogId() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
