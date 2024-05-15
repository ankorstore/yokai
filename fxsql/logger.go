package fxsql

import (
	"github.com/ankorstore/yokai/log"
)

// MigratorLogger is a logger compatible with [Goose].
//
// [Goose]: https://github.com/pressly/goose.
type MigratorLogger struct {
	logger *log.Logger
}

// NewMigratorLogger returns a new MigratorLogger instance.
func NewMigratorLogger(logger *log.Logger) *MigratorLogger {
	return &MigratorLogger{
		logger: logger,
	}
}

// Printf logs with info level.
func (l *MigratorLogger) Printf(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Fatalf logs with error level.
func (l *MigratorLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}
