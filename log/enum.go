package log

import (
	"strings"

	"github.com/rs/zerolog"
)

// FetchLogLevel returns a [Zerolog level] for a given value.
//
// [Zerolog level]: https://github.com/rs/zerolog/blob/master/log.go
//
//nolint:cyclop
func FetchLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "no-level":
		return zerolog.NoLevel
	case "disabled":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

// LogOutputWriter is an enum for the log output writers.
type LogOutputWriter int

const (
	StdoutOutputWriter LogOutputWriter = iota
	NoopOutputWriter
	TestOutputWriter
	ConsoleOutputWriter
)

// String returns a string representation of a [LogOutputWriter].
//
//nolint:exhaustive
func (l LogOutputWriter) String() string {
	switch l {
	case NoopOutputWriter:
		return Noop
	case TestOutputWriter:
		return Test
	case ConsoleOutputWriter:
		return Console
	default:
		return Stdout
	}
}

// FetchLogOutputWriter returns a [LogOutputWriter] for a given value.
func FetchLogOutputWriter(l string) LogOutputWriter {
	switch strings.ToLower(l) {
	case Noop:
		return NoopOutputWriter
	case Test:
		return TestOutputWriter
	case Console:
		return ConsoleOutputWriter
	default:
		return StdoutOutputWriter
	}
}
