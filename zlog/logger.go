package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

type LogConf struct {
	Level            string
	BufSwitch        bool
	BufSize          int
	BufFlushInterval time.Duration
	LogFilePath      string
}
type LogOption func(*LogConf)

func WithLevel(lv string) LogOption {
	return func(opt *LogConf) {
		opt.Level = lv
	}
}
func WithBufSwitch(sw bool) LogOption {
	return func(opt *LogConf) {
		opt.BufSwitch = sw
	}
}
func WithBufSize(buf int) LogOption {
	return func(opt *LogConf) {
		opt.BufSize = buf
	}
}
func WithBufFlushInterval(interval time.Duration) LogOption {
	return func(opt *LogConf) {
		opt.BufFlushInterval = interval
	}
}
func WithLogFilePath(LogFilePath string) LogOption {
	return func(opt *LogConf) {
		opt.LogFilePath = LogFilePath
	}
}

func defaultLogConf() *LogConf {
	return &LogConf{
		Level:            "INFO",
		BufSwitch:        true,
		BufSize:          128 * 1024, // 128kb
		BufFlushInterval: 3 * time.Second,
	}
}

func GetLogConfWithOpts(opts ...LogOption) LogConf {
	lc := defaultLogConf()
	for _, apply := range opts {
		apply(lc)
	}
	return *lc
}

func newLogger(lc LogConf) *zap.Logger {
	core := zapcore.NewCore(getLogEncoder(), getLogWriter(lc), getLogLevel(lc.Level))
	return zap.New(core)
}

func getLogEncoder() zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999"),
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter(conf LogConf) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if conf.LogFilePath != "" {
		w = &lumberjack.Logger{
			Filename:   conf.LogFilePath,
			MaxSize:    200,
			MaxBackups: 5,
			MaxAge:     14,   // days
			Compress:   true, // disabled by default
		}
	} else {
		w = os.Stdout
	}
	if !conf.BufSwitch {
		return zapcore.AddSync(w)
	}
	// 开启缓冲区
	return &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          conf.BufSize,
		FlushInterval: conf.BufFlushInterval,
	}
}

func getLogLevel(lv string) zapcore.Level {
	switch strings.ToUpper(lv) {
	case "DEBUG":
		return zap.DebugLevel
	case "INFO":
		return zap.InfoLevel
	case "WARN":
		return zap.WarnLevel
	case "ERROR":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
