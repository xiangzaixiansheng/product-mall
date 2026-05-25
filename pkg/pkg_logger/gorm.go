package pkg_logger

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GORMLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewGORMLogger() *GORMLogger {
	return &GORMLogger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *GORMLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GORMLogger) Info(ctx context.Context, s string, args ...any) {
	GetLoggerWithCtx(ctx).Info(s, "args", args)
}

func (l *GORMLogger) Warn(ctx context.Context, s string, args ...any) {
	GetLoggerWithCtx(ctx).Warn(s, "args", args)
}

func (l *GORMLogger) Error(ctx context.Context, s string, args ...any) {
	GetLoggerWithCtx(ctx).Error(s, "args", args)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	logger := GetLoggerWithCtx(ctx)
	if l.SourceField != "" {
		logger = logger.With("source", utils.FileWithLineNum())
	}
	if err != nil && !(err == gorm.ErrRecordNotFound && l.SkipErrRecordNotFound) {
		logger.Error("sql error", "sql", sql, "elapsed", elapsed, "error", err)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger.Warn("slow sql", "sql", sql, "elapsed", elapsed)
		return
	}

	if l.Debug {
		logger.Debug("sql trace", "sql", sql, "elapsed", elapsed)
		return
	}
}

func GetLoggerWithCtx(ctx context.Context) *slog.Logger {
	requestId := ""
	if ctx != nil {
		if id, ok := ctx.Value("X-Request-ID").(string); ok {
			requestId = id
		}
	}
	return LogrusObj.With("X-Request-ID", requestId)
}