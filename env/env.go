package env

import (
	"github.com/gin-gonic/gin"
	"os"
)

const (
	DefaultRootPath = "."
	RunEnv          = "RUN_ENV"

	RunEnvDev  = "dev"
	RunEnvTest = "test"
	RunEnvProd = "prod"
)

var (
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
	gin.SetMode(RunMode)
}

func SetRootPath(r string) {
	rootPath = r
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
