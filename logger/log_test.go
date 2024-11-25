package logger

import (
	"context"
	"github.com/pkg/errors"
	"testing"
)

func TestLog(t *testing.T) {
	Setup()

	ctx := WithTraceCtx(context.Background())
	Infof(ctx, "print: %s", "aaa")
	Warn(ctx, "bbb")
	Error(ctx, "ccc")

	WithCtx(ctx).Info().Msgf("print: %s", "test")

	ctx2 := WithTraceCtx(context.WithValue(context.Background(), CtxTraceID, "5678"))
	Info(ctx2, "ddd")

	Err(context.Background(), errors.New("err msg"), "fff")
}
