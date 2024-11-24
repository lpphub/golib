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
		Filename string
	}

	LogOption func(*LogConf)
)

var (
	logger *Logger
	once   sync.Once
)

func Setup(opts ...LogOption) {
	lc := defaultLogConf()
	for _, apply := range opts {
		apply(lc)
	}

	once.Do(func() {
		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if lc.Filename != "" {
			fileLogger := &lumberjack.Logger{
				Filename:   lc.Filename,
				MaxSize:    200,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   false,
			}
			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339
		log := zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
		logger = &log
	})
}

func defaultLogConf() *LogConf {
	return &LogConf{
		//Filename: "app.log",
	}
}

func Get() *Logger {
	return logger
}

func WithFilename(filename string) LogOption {
	return func(lc *LogConf) {
		lc.Filename = filename
	}
}
