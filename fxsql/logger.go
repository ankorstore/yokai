package fxsql

import (
	"fmt"

	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

// MigratorLogger is a logger compatible with [Goose].
//
// [Goose]: https://github.com/pressly/goose.
type MigratorLogger struct {
	logger *log.Logger
	stdout bool
}

// NewMigratorLogger returns a new MigratorLogger instance.
func NewMigratorLogger(logger *log.Logger, stdout bool) *MigratorLogger {
	return &MigratorLogger{
		logger: logger,
		stdout: stdout,
	}
}

// Printf logs with info level, and prints to stdout if configured to do so.
func (l *MigratorLogger) Printf(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)

	if l.stdout {
		//nolint:forbidigo
		fmt.Printf(format, v...)
	}
}

// Fatalf logs with fatal level, and prints to stdout if configured to do so.
func (l *MigratorLogger) Fatalf(format string, v ...interface{}) {
	l.logger.WithLevel(zerolog.FatalLevel).Msgf(format, v...)

	if l.stdout {
		//nolint:forbidigo
		fmt.Printf("[FATAL] "+format, v...)
	}
}
