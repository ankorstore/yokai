package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

// CtxOrmLogger is a logger compatible with the [Gorm logger].
//
// [Gorm logger]: https://gorm.io/docs/logger.html
type CtxOrmLogger struct {
	level      logger.LogLevel
	withValues bool
}

// NewCtxOrmLogger returns a new [CtxOrmLogger].
func NewCtxOrmLogger(level logger.LogLevel, withValues bool) *CtxOrmLogger {
	return &CtxOrmLogger{
		level:      level,
		withValues: withValues,
	}
}

// LogMode sets the logger log level.
func (l *CtxOrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level

	return l
}

// Error logs with error level.
func (l *CtxOrmLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Error {
		log.CtxLogger(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
	}
}

// Warn logs with warn level.
func (l *CtxOrmLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Warn {
		log.CtxLogger(ctx).Warn().Msg(fmt.Sprintf(msg, opts...))
	}
}

// Info logs with info level.
func (l *CtxOrmLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Info {
		log.CtxLogger(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
	}
}

// Trace logs with trace level.
func (l *CtxOrmLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	var event *zerolog.Event
	if err != nil {
		event = log.CtxLogger(ctx).Error()
	} else {
		event = log.CtxLogger(ctx).Debug()
	}

	event.Str("latency", time.Since(begin).String())

	sql, rows := f()
	if sql != "" {
		event.Str("sqlQuery", sql)
	}
	if rows > -1 {
		event.Int64("sqlRows", rows)
	}

	event.Send()
}

// ParamsFilter is used to filter SQL queries params from logging (replace with ?).
func (l *CtxOrmLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if !l.withValues {
		return sql, nil
	}

	return sql, params
}
