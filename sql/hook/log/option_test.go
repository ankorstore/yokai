package log_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestWithLevel(t *testing.T) {
	t.Parallel()

	opt := log.DefaultLogHookOptions()
	log.WithLevel(zerolog.DebugLevel)(&opt)

	assert.Equal(t, zerolog.DebugLevel, opt.Level)
}

func TestWithArguments(t *testing.T) {
	t.Parallel()

	opt := log.DefaultLogHookOptions()
	log.WithArguments(true)(&opt)

	assert.Equal(t, true, opt.Arguments)
}

func TestWithExcludedOperations(t *testing.T) {
	t.Parallel()

	exclusions := []sql.Operation{
		sql.ConnectionPingOperation,
		sql.ConnectionResetSessionOperation,
	}

	opt := log.DefaultLogHookOptions()
	log.WithExcludedOperations(exclusions...)(&opt)

	assert.Equal(t, exclusions, opt.ExcludedOperations)
}
