package logger

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	TraceID = "ctx_traceId"
)

func WithCtx(ctx context.Context) context.Context {
	traceID := ""
	if tid := ctx.Value(TraceID); tid != nil {
		traceID = tid.(string)
	} else {
		traceID = GenerateTraceID()
	}

	log := logger.With().Str("traceID", traceID).Logger()
	return log.WithContext(ctx)
}

func FromCtx(ctx context.Context) *Logger {
	if ctx == nil {
		return logger
	}
	if l := zerolog.Ctx(ctx); l != nil && l.GetLevel() != zerolog.Disabled {
		return l
	}
	return logger
}

func Info(ctx context.Context, msg string) {
	FromCtx(ctx).Info().Caller(1).Msg(msg)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Info().Caller(1).Msgf(format, v...)
}

func Error(ctx context.Context, msg string) {
	FromCtx(ctx).Error().Caller(1).Msg(msg)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Error().Caller(1).Msgf(format, v...)
}

func Err(ctx context.Context, err error, msg string) {
	FromCtx(ctx).Err(err).Caller(1).Msg(msg)
}

func Debug(ctx context.Context, msg string) {
	FromCtx(ctx).Debug().Caller(1).Msg(msg)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Debug().Caller(1).Msgf(format, v...)
}

func Warn(ctx context.Context, msg string) {
	FromCtx(ctx).Warn().Caller(1).Msg(msg)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Warn().Caller(1).Msgf(format, v...)
}

func Trace(ctx context.Context, msg string) {
	FromCtx(ctx).Trace().Caller(1).Msg(msg)
}

func Tracef(ctx context.Context, format string, v ...interface{}) {
	FromCtx(ctx).Trace().Caller(1).Msgf(format, v...)
}

func GenerateTraceID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
