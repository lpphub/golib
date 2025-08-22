package test

import (
	"context"
	"testing"

	"github.com/lpphub/golib/logger"
	"github.com/pkg/errors"
)

func TestPrint(t *testing.T) {
	logger.Setup()

	ctx := logger.WithCtx(context.Background())
	logger.Infof(ctx, "print: %s", "aaa")
	logger.Info(ctx, "bbb")
	logger.Error(ctx, "ccc")

	logger.FromCtx(ctx).Info().Msgf("print: %s", "test")

	ctx2 := logger.WithCtx(context.WithValue(context.Background(), logger.LogId, "5678"))
	logger.Info(ctx2, "ddd")

	logger.Err(context.Background(), errors.New("err msg"), "fff")
}
