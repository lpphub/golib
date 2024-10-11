package ware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BootstrapConf struct {
	TraceLog       TraceLogConfig
	CustomRecovery gin.RecoveryFunc
	Cors           bool
}

func Bootstraps(app *gin.Engine, opt BootstrapConf) {
	gin.SetMode(env.RunMode)

	// 中间件
	if opt.TraceLog.Enable {
		app.Use(TraceLog(opt.TraceLog))
	}
	app.Use(gin.CustomRecovery(opt.CustomRecovery))
	if opt.Cors {
		app.Use(Cors())
	}

	// ready check
	app.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"STATUS": "UP"})
	})
}

func ListenAndServe(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen to serve err: %s\n", err.Error())
		}
	}()
	log.Printf("Listening and serving HTTP on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err.Error())
	}
	log.Printf("Server exited")
}
