package hooktest_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/hooktest"
	"github.com/stretchr/testify/assert"
)

func TestWithSystem(t *testing.T) {
	t.Parallel()

	opt := hooktest.DefaultTestHookEventOptions()
	hooktest.WithSystem(sql.MysqlSystem)(&opt)

	assert.Equal(t, sql.MysqlSystem, opt.System)
}

func TestWithOperation(t *testing.T) {
	t.Parallel()

	opt := hooktest.DefaultTestHookEventOptions()
	hooktest.WithOperation(sql.ConnectionPingOperation)(&opt)

	assert.Equal(t, sql.ConnectionPingOperation, opt.Operation)
}

func TestWithQuery(t *testing.T) {
	t.Parallel()

	opt := hooktest.DefaultTestHookEventOptions()
	hooktest.WithQuery("SELECT * FROM foo")(&opt)

	assert.Equal(t, "SELECT * FROM foo", opt.Query)
}

func TestWithArguments(t *testing.T) {
	t.Parallel()

	opt := hooktest.DefaultTestHookEventOptions()
	hooktest.WithArguments(42)(&opt)

	assert.Equal(t, 42, opt.Arguments)
}
