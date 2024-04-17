package httpserver

import (
	"fmt"
	"io"
	"sync"

	"github.com/ankorstore/yokai/log"
	echologger "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// EchoLogger is a [log.Logger] wrapper for [echo.Logger] compatibility.
type EchoLogger struct {
	logger *log.Logger
	mutex  sync.RWMutex
	prefix string
}

// NewEchoLogger returns a new [log.Logger].
func NewEchoLogger(logger *log.Logger) *EchoLogger {
	return &EchoLogger{
		logger: logger,
	}
}

// ToZerolog converts new [log.Logger] to a [zerolog.Logger].
func (e *EchoLogger) ToZerolog() *zerolog.Logger {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.logger.ToZerolog()
}

// Output returns the output writer.
func (e *EchoLogger) Output() io.Writer {
	return e.logger
}

// SetOutput sets the output writer.
func (e *EchoLogger) SetOutput(w io.Writer) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.logger.Output(w)
}

// Prefix returns the log prefix.
func (e *EchoLogger) Prefix() string {
	return e.prefix
}

// SetPrefix sets the log prefix.
func (e *EchoLogger) SetPrefix(p string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.prefix = p
}

// Level returns the log level.
func (e *EchoLogger) Level() echologger.Lvl {
	return convertZeroLevel(e.logger.GetLevel())
}

// SetLevel sets the log level.
func (e *EchoLogger) SetLevel(v echologger.Lvl) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.logger = log.FromZerolog(e.logger.ToZerolog().Level(convertEchoLevel(v)))
}

// SetHeader sets the log header field.
func (e *EchoLogger) SetHeader(h string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.logger = log.FromZerolog(e.logger.With().Str("header", h).Logger())
}

// Debug produces a log with debug level.
func (e *EchoLogger) Debug(i ...interface{}) {
	e.logger.Debug().Msg(fmt.Sprint(i...))
}

// Debugf produces a formatted log with debug level.
func (e *EchoLogger) Debugf(format string, args ...interface{}) {
	e.logger.Debug().Msgf(format, args...)
}

// Debugj produces a json log with debug level.
func (e *EchoLogger) Debugj(j echologger.JSON) {
	e.logJSON(e.logger.Debug(), j)
}

// Info produces a log with info level.
func (e *EchoLogger) Info(i ...interface{}) {
	e.logger.Info().Msg(fmt.Sprint(i...))
}

// Infof produces a formatted log with info level.
func (e *EchoLogger) Infof(format string, args ...interface{}) {
	e.logger.Info().Msgf(format, args...)
}

// Infoj produces a json log with info level.
func (e *EchoLogger) Infoj(j echologger.JSON) {
	e.logJSON(e.logger.Info(), j)
}

// Warn produces a log with warn level.
func (e *EchoLogger) Warn(i ...interface{}) {
	e.logger.Warn().Msg(fmt.Sprint(i...))
}

// Warnf produces a formatted log with warn level.
func (e *EchoLogger) Warnf(format string, args ...interface{}) {
	e.logger.Warn().Msgf(format, args...)
}

// Warnj produces a json log with warn level.
func (e *EchoLogger) Warnj(j echologger.JSON) {
	e.logJSON(e.logger.Warn(), j)
}

// Error produces a log with error level.
func (e *EchoLogger) Error(i ...interface{}) {
	e.logger.Error().Msg(fmt.Sprint(i...))
}

// Errorf produces a formatted log with error level.
func (e *EchoLogger) Errorf(format string, args ...interface{}) {
	e.logger.Error().Msgf(format, args...)
}

// Errorj produces a json log with error level.
func (e *EchoLogger) Errorj(j echologger.JSON) {
	e.logJSON(e.logger.Error(), j)
}

// Fatal produces a log with fatal level.
func (e *EchoLogger) Fatal(i ...interface{}) {
	e.logger.Fatal().Msg(fmt.Sprint(i...))
}

// Fatalf produces a formatted log with fatal level.
func (e *EchoLogger) Fatalf(format string, args ...interface{}) {
	e.logger.Fatal().Msgf(format, args...)
}

// Fatalj produces a json log with fatal level.
func (e *EchoLogger) Fatalj(j echologger.JSON) {
	e.logJSON(e.logger.Fatal(), j)
}

// Panic produces a log with panic level.
func (e *EchoLogger) Panic(i ...interface{}) {
	e.logger.Panic().Msg(fmt.Sprint(i...))
}

// Panicf produces a formatted log with panic level.
func (e *EchoLogger) Panicf(format string, args ...interface{}) {
	e.logger.Panic().Msgf(format, args...)
}

// Panicj produces a json log with panic level.
func (e *EchoLogger) Panicj(j echologger.JSON) {
	e.logJSON(e.logger.Panic(), j)
}

// Print produces a log with no level.
func (e *EchoLogger) Print(i ...interface{}) {
	e.logger.WithLevel(zerolog.NoLevel).Str("level", "---").Msg(fmt.Sprint(i...))
}

// Printf produces a formatted log with no level.
func (e *EchoLogger) Printf(format string, i ...interface{}) {
	e.logger.WithLevel(zerolog.NoLevel).Str("level", "---").Msgf(format, i...)
}

// Printj produces a json log with no level.
func (e *EchoLogger) Printj(j echologger.JSON) {
	e.logJSON(e.logger.WithLevel(zerolog.NoLevel).Str("level", "---"), j)
}

func (e *EchoLogger) logJSON(event *zerolog.Event, j echologger.JSON) {
	for k, v := range j {
		event = event.Interface(k, v)
	}

	event.Msg("")
}

func convertZeroLevel(level zerolog.Level) echologger.Lvl {
	switch level {
	case zerolog.TraceLevel, zerolog.DebugLevel:
		return echologger.DEBUG
	case zerolog.InfoLevel:
		return echologger.INFO
	case zerolog.WarnLevel:
		return echologger.WARN
	case zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel:
		return echologger.ERROR
	case zerolog.NoLevel, zerolog.Disabled:
		return echologger.OFF
	default:
		return echologger.INFO
	}
}

func convertEchoLevel(level echologger.Lvl) zerolog.Level {
	switch level {
	case echologger.DEBUG:
		return zerolog.DebugLevel
	case echologger.INFO:
		return zerolog.InfoLevel
	case echologger.WARN:
		return zerolog.WarnLevel
	case echologger.ERROR:
		return zerolog.ErrorLevel
	case echologger.OFF:
		return zerolog.NoLevel
	default:
		return zerolog.InfoLevel
	}
}
