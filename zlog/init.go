package zlog

import (
	"go.uber.org/zap"
)

var (
	ZapLogger     *zap.Logger
	SugaredLogger *zap.SugaredLogger

	logConf LogConf
)

func InitLog(opts ...LogOption) {
	logConf = GetLogConfWithOpts(opts...)
	if SugaredLogger == nil {
		if ZapLogger == nil {
			ZapLogger = newLogger()
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
