package hooktest

import "github.com/ankorstore/yokai/sql"

var (
	TestHookEventQuery    = "SELECT * FROM foo WHERE id = ?"
	TestHookEventArgument = "42"
)

// Options are the options for NewTestHookEvent.
type Options struct {
	System    sql.System
	Operation sql.Operation
	Query     string
	Arguments any
}

// DefaultTestHookEventOptions are the default options for NewTestHookEvent.
func DefaultTestHookEventOptions() Options {
	return Options{
		System:    sql.SqliteSystem,
		Operation: sql.ConnectionQueryOperation,
		Query:     TestHookEventQuery,
		Arguments: TestHookEventArgument,
	}
}

// TestHookEventOption are the functional options for NewTestHookEvent.
type TestHookEventOption func(o *Options)

// WithSystem is used to set the sql.System of the test sql.HookEvent.
func WithSystem(system sql.System) TestHookEventOption {
	return func(o *Options) {
		o.System = system
	}
}

// WithOperation is used to set the sql.Operation of the test sql.HookEvent.
func WithOperation(operation sql.Operation) TestHookEventOption {
	return func(o *Options) {
		o.Operation = operation
	}
}

// WithQuery is used to set the query of the test sql.HookEvent.
func WithQuery(query string) TestHookEventOption {
	return func(o *Options) {
		o.Query = query
	}
}

// WithArguments is used to set the arguments of the test sql.HookEvent.
func WithArguments(arguments any) TestHookEventOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}
