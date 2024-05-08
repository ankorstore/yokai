package trace_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/stretchr/testify/assert"
)

func TestWithArguments(t *testing.T) {
	t.Parallel()

	opt := trace.DefaultTraceHookOptions()
	trace.WithArguments(true)(&opt)

	assert.Equal(t, true, opt.Arguments)
}

func TestWithExcludedOperations(t *testing.T) {
	t.Parallel()

	exclusions := []sql.Operation{
		sql.ConnectionPingOperation,
		sql.ConnectionResetSessionOperation,
	}

	opt := trace.DefaultTraceHookOptions()
	trace.WithExcludedOperations(exclusions...)(&opt)

	assert.Equal(t, exclusions, opt.ExcludedOperations)
}
