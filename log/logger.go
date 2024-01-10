package log

import (
	"github.com/rs/zerolog"
)

const (
	Level   = "level"
	Message = "message"
	Service = "service"
	Time    = "time"
	Stdout  = "stdout"
	Noop    = "noop"
	Test    = "test"
	Console = "console"
)

// Logger provides the possibility to generate logs, and inherits of all [Zerolog] features.
//
// [Zerolog]: https://github.com/rs/zerolog/tree/master
type Logger struct {
	*zerolog.Logger
}

// ToZerolog converts as [Logger] into a [Zerolog logger].
//
// [Zerolog logger]: https://github.com/rs/zerolog/blob/master/log.go
func (l *Logger) ToZerolog() *zerolog.Logger {
	return l.Logger
}

// FromZerolog converts as [Zerolog logger] into a [Logger].
//
// [Zerolog logger]: https://github.com/rs/zerolog/blob/master/log.go
func FromZerolog(logger zerolog.Logger) *Logger {
	return &Logger{&logger}
}
