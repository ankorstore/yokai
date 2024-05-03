package hook_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql/hook"
	"github.com/stretchr/testify/assert"
)

func TestNewHookEvent(t *testing.T) {
	t.Parallel()

	event := createTestEvent()
	assert.IsType(t, &hook.HookEvent{}, event)

	assert.Equal(t, "system", event.System())
	assert.Equal(t, "operation", event.Operation())
	assert.Equal(t, "query", event.Query())
	assert.Equal(t, "argument", event.Arguments())
}

func TestHookEventLastInsertId(t *testing.T) {
	t.Parallel()

	event := createTestEvent()
	event.SetLastInsertId(int64(1))

	assert.Equal(t, int64(1), event.LastInsertId())
}

func TestHookEventRowsAffected(t *testing.T) {
	t.Parallel()

	event := createTestEvent()
	event.SetRowsAffected(int64(1))

	assert.Equal(t, int64(1), event.RowsAffected())
}

func TestHookEventError(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("test error")

	event := createTestEvent()
	event.SetError(err)

	assert.Equal(t, err, event.Error())
}

func TestHookEventLatency(t *testing.T) {
	t.Parallel()

	event := createTestEvent()

	_, err := event.Latency()
	assert.Error(t, err)
	assert.Equal(t, "event was not started", err.Error())

	event.Start()

	_, err = event.Latency()
	assert.Error(t, err)
	assert.Equal(t, "event was not stopped", err.Error())

	event.Stop()

	latency, err := event.Latency()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, latency.Nanoseconds(), int64(0))
}

func createTestEvent() *hook.HookEvent {
	return hook.NewHookEvent("system", "operation", "query", "argument")
}
