package log

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/sql"
)

// LogHook is a hook.Hook implementation for SQL logging.
type LogHook struct {
	options Options
}

// NewLogHook returns a new LogHook, for a provided list of LogHookOption.
func NewLogHook(options ...LogHookOption) *LogHook {
	appliedOpts := DefaultLogHookOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	return &LogHook{
		options: appliedOpts,
	}
}

// Before executes SQL logging logic before SQL operations.
func (h *LogHook) Before(ctx context.Context, _ *sql.HookEvent) context.Context {
	return ctx
}

// After executes SQL logging logic after SQL operations.
func (h *LogHook) After(ctx context.Context, event *sql.HookEvent) {
	if sql.ContainsOperation(h.options.ExcludedOperations, event.Operation()) {
		return
	}

	logger := log.CtxLogger(ctx)

	loggerEvent := logger.WithLevel(h.options.Level)
	if event.Error() != nil {
		if !errors.Is(event.Error(), driver.ErrSkip) {
			loggerEvent = logger.Error().Err(event.Error())
		}
	}

	loggerEvent.
		Str("system", event.System().String()).
		Str("operation", event.Operation().String()).
		Int64("lastInsertId", event.LastInsertId()).
		Int64("rowsAffected", event.RowsAffected())

	if event.Query() != "" {
		loggerEvent.Str("query", event.Query())
	}

	if h.options.Arguments && event.Arguments() != nil {
		loggerEvent.Interface("arguments", event.Arguments())
	}

	latency, err := event.Latency()
	if err == nil {
		loggerEvent.Str("latency", latency.String())
	}

	loggerEvent.Msg("sql logger")
}
