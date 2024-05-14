package ware

import (
	"github.com/gin-gonic/gin"
)

type BootstrapConf struct {
	OpenTrace      bool
	Cors           bool
	CustomRecovery gin.RecoveryFunc
}

func Bootstrap(app *gin.Engine, opt BootstrapConf) {
	if opt.OpenTrace {
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
