package model

import (
	"context"
	"time"

	"log/slog"

	"gorm.io/gorm/logger"
)

// tracelogger 集成traceid
type tracelogger struct {
	lvl logger.LogLevel
}

func (l *tracelogger) LogMode(lvl logger.LogLevel) logger.Interface {
	newlog := *l
	newlog.lvl = lvl
	return &newlog
}
func (l *tracelogger) Info(ctx context.Context, s string, v ...interface{})  {}
func (l *tracelogger) Warn(ctx context.Context, s string, v ...interface{})  {}
func (l *tracelogger) Error(ctx context.Context, s string, v ...interface{}) {}
func (l *tracelogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.lvl <= logger.Silent {
		return
	}
	sql, rows := fc()
	dur := time.Since(begin)
	logattr := []any{
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Duration("latency", dur),
	}
	switch {
	case err != nil && l.lvl >= logger.Error:
		slog.ErrorContext(ctx, "gorm.trace", append(logattr, slog.String("err", err.Error()))...)
	case dur >= time.Millisecond*500 && l.lvl >= logger.Warn:
		slog.WarnContext(ctx, "gorm.trace slow sql", logattr...)
	case l.lvl == logger.Info:
		slog.InfoContext(ctx, "gorm.trace", logattr...)
	}
}
