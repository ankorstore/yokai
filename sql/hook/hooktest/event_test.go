package hooktest_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/hooktest"
	"github.com/stretchr/testify/assert"
)

func TestNewTestHookEventWithDefaults(t *testing.T) {
	t.Parallel()

	event := hooktest.NewTestHookEvent()

	assert.Equal(t, sql.SqliteSystem, event.System())
	assert.Equal(t, sql.ConnectionQueryOperation, event.Operation())
	assert.Equal(t, hooktest.TestHookEventQuery, event.Query())
	assert.Equal(t, hooktest.TestHookEventArgument, event.Arguments())
}

func TestNewTestHookEventWithOptions(t *testing.T) {
	t.Parallel()

	event := hooktest.NewTestHookEvent(
		hooktest.WithSystem(sql.MysqlSystem),
		hooktest.WithOperation(sql.ConnectionPingOperation),
		hooktest.WithQuery("SELECT * FROM bar WHERE id = ?"),
		hooktest.WithArguments(24),
	)

	assert.Equal(t, sql.MysqlSystem, event.System())
	assert.Equal(t, sql.ConnectionPingOperation, event.Operation())
	assert.Equal(t, "SELECT * FROM bar WHERE id = ?", event.Query())
	assert.Equal(t, 24, event.Arguments())
}
