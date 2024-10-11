package ware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BootstrapConf struct {
	LogTrace       bool
	Cors           bool
	CustomRecovery gin.RecoveryFunc
}

func Bootstrap(app *gin.Engine, opt BootstrapConf) {
	if opt.LogTrace {
		app.Use(LogTrace())
	}

	app.Use(gin.CustomRecovery(opt.CustomRecovery))

	if opt.Cors {
		app.Use(Cors())
	}

	app.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"STATUS": "UP"})
	})
}

func ListenAndServe(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen to serve err: %s\n", err)
		}
	}()
	log.Printf("Listening and serving HTTP on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exited")
}
