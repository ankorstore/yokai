package trace

import "github.com/ankorstore/yokai/sql"

// Options are the options for TraceHook.
type Options struct {
	Arguments          bool
	ExcludedOperations []sql.Operation
}

// DefaultTraceHookOptions are the default options for TraceHook.
func DefaultTraceHookOptions() Options {
	return Options{
		Arguments:          false,
		ExcludedOperations: []sql.Operation{},
	}
}

// TraceHookOption are the functional options for TraceHook.
type TraceHookOption func(o *Options)

// WithArguments is used to enable the SQL arguments tracing.
func WithArguments(arguments bool) TraceHookOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}

// WithExcludedOperations is used to exclude a list of database operations from tracing.
func WithExcludedOperations(excludedOperations ...sql.Operation) TraceHookOption {
	return func(o *Options) {
		o.ExcludedOperations = excludedOperations
	}
}
