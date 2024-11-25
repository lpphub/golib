package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type (
	Logger = zerolog.Logger

	LogConf struct {
		Stdout  bool
		LogFile string
	}
	LogOption func(*LogConf)
)

var (
	logger *Logger
	once   sync.Once
)

func Setup(opts ...LogOption) {
	once.Do(func() {
		lc := defaultLogConf()
		for _, apply := range opts {
			apply(lc)
		}

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		writers := make([]io.Writer, 0)
		if lc.Stdout {
			writers = append(writers, os.Stdout)
		}
		if lc.LogFile != "" {
			fileLogger := &lumberjack.Logger{
				Filename:   lc.LogFile,
				MaxSize:    200,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   false,
			}
			writers = append(writers, fileLogger)
		}
		if len(writers) > 0 {
			output = zerolog.MultiLevelWriter(writers...)
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339
		log := zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger()
		logger = &log
	})
}

func defaultLogConf() *LogConf {
	return &LogConf{
		//Stdout: true,
	}
}

func Get() *Logger {
	return logger
}

func WithLogFile(logFile string) LogOption {
	return func(lc *LogConf) {
		lc.LogFile = logFile
	}
}

func WithStdout(stdout bool) LogOption {
	return func(lc *LogConf) {
		lc.Stdout = stdout
	}
}
