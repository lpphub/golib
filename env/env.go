package env

import (
	"os"

	"github.com/gin-gonic/gin"
)

const (
	DefaultRootPath = "."
	RunEnv          = "RUN_ENV"

	RunEnvDev  = "dev"
	RunEnvTest = "test"
	RunEnvProd = "prod"
)

var (
	AppName string
	RunMode string

	runEnv   string
	rootPath string
)

func init() {
	RunMode = gin.ReleaseMode
	r := os.Getenv(RunEnv)
	switch r {
	case RunEnvProd:
		runEnv = RunEnvProd
	case RunEnvTest:
		runEnv = RunEnvTest
	default:
		runEnv = RunEnvDev
		RunMode = gin.DebugMode
	}
}

func SetRootPath(r string) {
	rootPath = r
}

func SetAppName(appName string) {
	AppName = appName
}

func GetRunEnv() string {
	return runEnv
}

func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}
