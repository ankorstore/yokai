package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Options are options for the [LoggerFactory] implementations.
type Options struct {
	ServiceName  string
	Level        zerolog.Level
	OutputWriter io.Writer
}

// DefaultLoggerOptions are the default options used in the [DefaultLoggerFactory].
func DefaultLoggerOptions() Options {
	return Options{
		ServiceName:  "default",
		Level:        zerolog.InfoLevel,
		OutputWriter: os.Stdout,
	}
}

// LoggerOption are functional options for the [LoggerFactory] implementations.
type LoggerOption func(o *Options)

// WithServiceName is used to add automatically a service log field value.
func WithServiceName(n string) LoggerOption {
	return func(o *Options) {
		o.ServiceName = n
	}
}

// WithLevel is used to specify the log level to use.
func WithLevel(l zerolog.Level) LoggerOption {
	return func(o *Options) {
		o.Level = l
	}
}

// WithOutputWriter is used to specify the output writer to use.
func WithOutputWriter(w io.Writer) LoggerOption {
	return func(o *Options) {
		o.OutputWriter = w
	}
}
