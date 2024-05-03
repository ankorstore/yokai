package log

import (
	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

// Options are the options for LogHook.
type Options struct {
	Level              zerolog.Level
	Arguments          bool
	ExcludedOperations []string
}

// DefaultLogHookOptions are the default options for LogHook.
func DefaultLogHookOptions() Options {
	return Options{
		Level:              zerolog.DebugLevel,
		Arguments:          false,
		ExcludedOperations: []string{},
	}
}

// LogHookOption are the functional options for LogHook.
type LogHookOption func(o *Options)

// WithLevel is used to configure the SQL logging level.
func WithLevel(level string) LogHookOption {
	return func(o *Options) {
		o.Level = log.FetchLogLevel(level)
	}
}

// WithArguments is used to enable the SQL arguments logging.
func WithArguments(arguments bool) LogHookOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}

// WithExcludedOperations is used to exclude a list of SQL operations from logging.
func WithExcludedOperations(excludedOperations ...string) LogHookOption {
	return func(o *Options) {
		o.ExcludedOperations = excludedOperations
	}
}
