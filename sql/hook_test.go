package sql_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/sqltest"
	"github.com/stretchr/testify/assert"
)

func TestNewHookEvent(t *testing.T) {
	t.Parallel()

	event := sqltest.NewTestHookEvent()
	assert.IsType(t, &sql.HookEvent{}, event)

	assert.Equal(t, sql.SqliteSystem, event.System())
	assert.Equal(t, sql.ConnectionQueryOperation, event.Operation())
	assert.Equal(t, sqltest.TestHookEventQuery, event.Query())
	assert.Equal(t, sqltest.TestHookEventArgument, event.Arguments())
}

func TestHookEventLastInsertId(t *testing.T) {
	t.Parallel()

	event := sqltest.NewTestHookEvent()
	event.SetLastInsertId(int64(1))

	assert.Equal(t, int64(1), event.LastInsertId())
}

func TestHookEventRowsAffected(t *testing.T) {
	t.Parallel()

	event := sqltest.NewTestHookEvent()
	event.SetRowsAffected(int64(1))

	assert.Equal(t, int64(1), event.RowsAffected())
}

func TestHookEventError(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("test error")

	event := sqltest.NewTestHookEvent()
	event.SetError(err)

	assert.Equal(t, err, event.Error())
}

func TestHookEventLatency(t *testing.T) {
	t.Parallel()

	event := sqltest.NewTestHookEvent()

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
