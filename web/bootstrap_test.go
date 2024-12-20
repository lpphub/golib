package web

import (
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger"
	"github.com/lpphub/golib/logger/logx"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	logger.Setup()

	r := gin.New()
	Bootstraps(r, BootstrapConf{
		Cors: true,
		AccessLog: AccessLogConfig{
			Enable:    true,
			SkipPaths: []string{"/metrics"},
		},
		CustomRecovery: func(ctx *gin.Context, err any) {
			JsonWithError(ctx, Error{-1, "test"})
		},
	})

	r.GET("/test", func(ctx *gin.Context) {
		logx.Infof(ctx, "哈哈: %s", "bb")

		logx.Err(ctx, errors.New("test"), "")

		JsonWithSuccess(ctx, "test")
	})

	ListenAndServe(&http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	})
}
