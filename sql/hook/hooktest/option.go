package hooktest

import "github.com/ankorstore/yokai/sql"

var (
	TestHookEventQuery    = "SELECT * FROM foo WHERE id = ?"
	TestHookEventArgument = "42"
)

// Options are the options for TraceHook.
type Options struct {
	System    sql.System
	Operation sql.Operation
	Query     string
	Arguments any
}

// DefaultTestHookEventOptions are the default options for TraceHook.
func DefaultTestHookEventOptions() Options {
	return Options{
		System:    sql.SqliteSystem,
		Operation: sql.ConnectionQueryOperation,
		Query:     TestHookEventQuery,
		Arguments: TestHookEventArgument,
	}
}

// TestHookEventOption are the functional options for TraceHook.
type TestHookEventOption func(o *Options)

// WithSystem is used to set the .
func WithSystem(system sql.System) TestHookEventOption {
	return func(o *Options) {
		o.System = system
	}
}

func WithOperation(operation sql.Operation) TestHookEventOption {
	return func(o *Options) {
		o.Operation = operation
	}
}

func WithQuery(query string) TestHookEventOption {
	return func(o *Options) {
		o.Query = query
	}
}

func WithArguments(arguments any) TestHookEventOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}
