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

	if SugaredLogger == nil {
		if ZapLogger == nil {
			ZapLogger = newLogger(lc)
		}
		SugaredLogger = ZapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	}
}

func Close() {
	if ZapLogger != nil {
		_ = ZapLogger.Sync()
	}
	if SugaredLogger != nil {
		_ = SugaredLogger.Sync()
	}
}
