package sql_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	hooks := []sql.Hook{
		trace.NewTraceHook(),
		log.NewLogHook(),
	}

	config := sql.NewConfiguration(sql.SqliteSystem, hooks...)

	assert.Equal(t, sql.SqliteSystem, config.System())
	assert.Equal(t, hooks, config.Hooks())
}
