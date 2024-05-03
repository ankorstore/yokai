package log_test

import (
	"context"
	"fmt"
	"testing"

	yokailog "github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/sql/hook"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogHookWithDefaults(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := yokailog.NewDefaultLoggerFactory().Create(
		yokailog.WithLevel(zerolog.DebugLevel),
		yokailog.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	h := log.NewLogHook()

	ctx := logger.WithContext(context.Background())
	event := hook.NewHookEvent("system", "operation", "query", "argument")

	newCtx := h.Before(ctx, event)
	assert.Same(t, ctx, newCtx)

	event.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(newCtx, event)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "debug",
		"system":       "system",
		"operation":    "operation",
		"query":        "query",
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})
}

func TestLogHookWithOptions(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := yokailog.NewDefaultLoggerFactory().Create(
		yokailog.WithLevel(zerolog.DebugLevel),
		yokailog.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	h := log.NewLogHook(
		log.WithLevel("info"),
		log.WithArguments(true),
		log.WithExcludedOperations("excludedOperation"),
	)

	ctx := logger.WithContext(context.Background())

	// regular event
	event := hook.NewHookEvent("system", "regularOperation", "query", "argument")

	newCtx := h.Before(ctx, event)
	assert.Same(t, ctx, newCtx)

	event.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(newCtx, event)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "info",
		"system":       "system",
		"operation":    "regularOperation",
		"query":        "query",
		"arguments":    "argument",
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})

	// excluded operation event
	excludedOperationEvent := hook.NewHookEvent("system", "excludedOperation", "query", "argument")

	h.Before(ctx, excludedOperationEvent)

	excludedOperationEvent.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(ctx, excludedOperationEvent)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "info",
		"system":       "system",
		"operation":    "excludedOperation",
		"query":        "query",
		"arguments":    "argument",
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})

	// error event
	errorEvent := hook.NewHookEvent("system", "errorOperation", "query", "argument")

	h.Before(ctx, errorEvent)

	errorEvent.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		SetError(fmt.Errorf("test error")).
		Stop()

	h.After(ctx, errorEvent)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "error",
		"error":        "test error",
		"system":       "system",
		"operation":    "errorOperation",
		"query":        "query",
		"arguments":    "argument",
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})
}
