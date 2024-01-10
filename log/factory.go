package log

import (
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

// LoggerFactory is the interface for [Logger] factories.
type LoggerFactory interface {
	Create(options ...LoggerOption) (*Logger, error)
}

// DefaultLoggerFactory is the default [LoggerFactory] implementation.
type DefaultLoggerFactory struct{}

// NewDefaultLoggerFactory returns a [DefaultLoggerFactory], implementing [LoggerFactory].
func NewDefaultLoggerFactory() LoggerFactory {
	return &DefaultLoggerFactory{}
}

// Create returns a new [Logger], and accepts a list of [LoggerOption].
// For example:
//
//	var logger, _ = log.NewDefaultLoggerFactory().Create()
//
// is equivalent to:
//
//	var logger, _ = log.NewDefaultLoggerFactory().Create(
//		log.WithServiceName("default"),   // adds {"service":"default"} to log records
//		log.WithLevel(zerolog.InfoLevel), // logs records with level >= info
//		log.WithOutputWriter(os.Stdout),  // sends logs records to stdout
//	)
func (f *DefaultLoggerFactory) Create(options ...LoggerOption) (*Logger, error) {
	appliedOpts := DefaultLoggerOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	logger := zerolog.
		New(appliedOpts.OutputWriter).
		With().
		Timestamp().
		Str(Service, appliedOpts.ServiceName).
		Logger().
		Level(appliedOpts.Level)

	return &Logger{&logger}, nil
}
