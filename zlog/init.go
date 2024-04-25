package zlog

import (
	"go.uber.org/zap"
)

var (
	ZapLogger     *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

func InitLog(opts ...LogOption) {
	lc := GetLogConfWithOpts(opts...)
	GetLogger(*lc)
}

func Close() {
	if ZapLogger != nil {
		_ = ZapLogger.Sync()
	}
	if SugaredLogger != nil {
		_ = SugaredLogger.Sync()
	}
}
