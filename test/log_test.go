package test

import (
	"context"
	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
	"testing"
)

func TestPrint(t *testing.T) {
	logger.Setup()

	ctx := logger.WithCtx(context.Background())
	logger.Infof(ctx, "print: %s", "aaa")
	logger.Warn(ctx, "bbb")
	logger.Error(ctx, "ccc")

	logger.FromCtx(ctx).Info().Msgf("print: %s", "test")

	ctx2 := logger.WithCtx(context.WithValue(context.Background(), logger.TraceID, "5678"))
	logger.Info(ctx2, "ddd")

	logger.Err(context.Background(), errors.New("err msg"), "fff")
}
