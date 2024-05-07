package log_test

import (
	"context"
	"fmt"
	"testing"

	yokailog "github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/hooktest"
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
	event := hooktest.NewTestHookEvent()

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
		"system":       event.System().String(),
		"operation":    event.Operation().String(),
		"query":        event.Query(),
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
		log.WithLevel(zerolog.InfoLevel),
		log.WithArguments(true),
		log.WithExcludedOperations(sql.ConnectionResetSessionOperation),
	)

	ctx := logger.WithContext(context.Background())

	// regular event
	event := hooktest.NewTestHookEvent()

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
		"system":       event.System().String(),
		"operation":    event.Operation().String(),
		"query":        event.Query(),
		"arguments":    event.Arguments(),
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})

	// excluded operation event
	logBuffer.Reset()

	excludedOperationEvent := hooktest.NewTestHookEvent(hooktest.WithOperation(sql.ConnectionResetSessionOperation))

	h.Before(ctx, excludedOperationEvent)

	excludedOperationEvent.
		Start().
		SetLastInsertId(1).
		SetRowsAffected(2).
		Stop()

	h.After(ctx, excludedOperationEvent)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "info",
		"system":       excludedOperationEvent.System().String(),
		"operation":    excludedOperationEvent.Operation().String(),
		"query":        event.Query(),
		"arguments":    event.Arguments(),
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})

	// error event
	logBuffer.Reset()

	errorEvent := hooktest.NewTestHookEvent()

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
		"system":       errorEvent.System().String(),
		"operation":    errorEvent.Operation().String(),
		"query":        event.Query(),
		"arguments":    event.Arguments(),
		"lastInsertId": "1",
		"rowsAffected": "2",
		"message":      "sql logger",
	})
}
