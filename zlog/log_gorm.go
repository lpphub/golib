package zlog

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	ormutil "gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	Addr      string
	Database  string
	MaxSqlLen int
	logger    *zap.Logger
}

func NewGormLogger(db, addr string, sqlLen int) logger.Interface {
	if ZapLogger == nil {
		return logger.Default
	}
	if sqlLen == 0 {
		sqlLen = 1024
	}
	return &GormLogger{
		Database:  db,
		Addr:      addr,
		MaxSqlLen: sqlLen,
		logger:    ZapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(3)),
	}
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	// 非trace日志改为debug级别输出
	l.logger.Debug(m, l.commonFields(ctx)...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	l.logger.Warn(m, l.commonFields(ctx)...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormutil.FileWithLineNum()}, data...)...)
	l.logger.Error(m, l.commonFields(ctx)...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Now().Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	// 请求是否成功
	msg := "mysql do success"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
	}

	sql, rows := fc()
	if l.MaxSqlLen <= 0 {
		sql = ""
	} else if len(sql) > l.MaxSqlLen {
		sql = sql[:l.MaxSqlLen]
	}

	fields := append(l.commonFields(ctx),
		zap.Int64("affectedRow", rows),
		//zap.String("fileLine", ormutil.FileWithLineNum()),
		zap.Float64("cost", cost),
		zap.String("sql", sql),
	)
	l.logger.Info(msg, fields...)
}

func (l *GormLogger) commonFields(ctx context.Context) []zap.Field {
	var logId string
	if c, ok := ctx.(*gin.Context); ok && c != nil {
		logId, _ = ctx.Value(GinCtxLogId).(string)
	}

	fields := []zap.Field{
		zap.String("logId", logId),
		zap.String("service", "mysql"),
		zap.String("db", l.Database),
		zap.String("addr", l.Addr),
	}
	return fields
}
