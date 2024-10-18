package zlog

import (
	"fmt"
	"github.com/lpphub/golib/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	_logInfo   = "info"
	_logError  = "error"
	_logStdout = "stdout"
)

type LogConf struct {
	LogLevel         string
	LogPath          string
	BufSwitch        bool
	BufSize          int
	BufFlushInterval time.Duration
}
type LogOption func(*LogConf)

func WithLogLevel(lv string) LogOption {
	return func(opt *LogConf) {
		opt.LogLevel = lv
	}
}
func WithLogPath(LogPath string) LogOption {
	return func(opt *LogConf) {
		opt.LogPath = LogPath
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

func defaultLogConf() *LogConf {
	return &LogConf{
		LogLevel:         "INFO",
		BufSwitch:        true,
		BufSize:          128 * 1024, // 128kb
		BufFlushInterval: 3 * time.Second,
	}
}

func GetLogConfWithOpts(opts ...LogOption) *LogConf {
	lc := defaultLogConf()
	for _, apply := range opts {
		apply(lc)
	}
	return lc
}

func newLogger() *zap.Logger {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl >= getLogLevel(logConf.LogLevel)
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= getLogLevel(logConf.LogLevel)
	})

	if logConf.LogPath != "" {
		core := zapcore.NewTee(
			zapcore.NewCore(getLogEncoder(), getLogWriter(_logInfo), infoLevel),
			zapcore.NewCore(getLogEncoder(), getLogWriter(_logError), errorLevel),
		)
		return zap.New(core)
	}

	// 控制台输出
	return zap.New(zapcore.NewCore(getLogEncoder(), getLogWriter(_logStdout), getLogLevel(logConf.LogLevel)))
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

func getLogWriter(logType string) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if logType == _logStdout {
		w = os.Stdout
	} else {
		app := env.AppName
		if app == "" {
			app = "server"
		}
		filename := filepath.Join(strings.TrimSuffix(logConf.LogPath, "/"), fmt.Sprintf("%s_%s.log", app, logType))
		w = &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    200,
			MaxBackups: 5,
			MaxAge:     14,    // days
			Compress:   false, // disabled by default
		}
	}
	if !logConf.BufSwitch {
		return zapcore.AddSync(w)
	}
	// 开启缓冲区
	return &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          logConf.BufSize,
		FlushInterval: logConf.BufFlushInterval,
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
