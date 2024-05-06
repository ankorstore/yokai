package log

import (
	"github.com/ankorstore/yokai/sql"
	"github.com/rs/zerolog"
)

// Options are the options for LogHook.
type Options struct {
	Level              zerolog.Level
	Arguments          bool
	ExcludedOperations []sql.Operation
}

// DefaultLogHookOptions are the default options for LogHook.
func DefaultLogHookOptions() Options {
	return Options{
		Level:              zerolog.DebugLevel,
		Arguments:          false,
		ExcludedOperations: []sql.Operation{},
	}
}

// LogHookOption are the functional options for LogHook.
type LogHookOption func(o *Options)

// WithLevel is used to configure the SQL logging level.
func WithLevel(level zerolog.Level) LogHookOption {
	return func(o *Options) {
		o.Level = level
	}
}

// WithArguments is used to enable the SQL arguments logging.
func WithArguments(arguments bool) LogHookOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}

// WithExcludedOperations is used to exclude a list of SQL operations from logging.
func WithExcludedOperations(excludedOperations ...sql.Operation) LogHookOption {
	return func(o *Options) {
		o.ExcludedOperations = excludedOperations
	}
}
