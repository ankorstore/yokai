package hooktest

import "github.com/ankorstore/yokai/sql"

// NewTestHookEvent returns a new sql.HookEvent for tests purposes.
func NewTestHookEvent(options ...TestHookEventOption) *sql.HookEvent {
	appliedOpts := DefaultTestHookEventOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	return sql.NewHookEvent(
		appliedOpts.System,
		appliedOpts.Operation,
		appliedOpts.Query,
		appliedOpts.Arguments,
	)
}
