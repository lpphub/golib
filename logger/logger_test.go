package logger

import (
	"context"
	"github.com/pkg/errors"
	"testing"
)

func TestGet(t *testing.T) {
	Setup()

	ctx := WithCtx(context.Background())
	Infof(ctx, "print: %s", "aaa")
	Warn(ctx, "bbb")
	Error(ctx, "ccc")

	FromCtx(ctx).Info().Msgf("print: %s", "test")

	ctx2 := WithCtx(context.WithValue(context.Background(), CtxTraceID, "5678"))
	Info(ctx2, "ddd")

	Err(ctx2, errors.New("err msg"), "fff")
}