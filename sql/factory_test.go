package sql_test

import (
	"database/sql/driver"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultDriverFactory(t *testing.T) {
	t.Parallel()

	factory := sql.NewDefaultDriverFactory()

	assert.IsType(t, &sql.DefaultDriverFactory{}, factory)
	assert.Implements(t, (*sql.DriverFactory)(nil), factory)
}

func TestDefaultDriverFactoryCreate(t *testing.T) {
	t.Parallel()

	factory := sql.NewDefaultDriverFactory()

	tests := []struct {
		system        sql.System
		hooks         []sql.Hook
		expectedBase  driver.Driver
		expectedError bool
	}{
		{
			sql.SqliteSystem,
			[]sql.Hook{log.NewLogHook()},
			&sqlite3.SQLiteDriver{},
			false,
		},
		{
			sql.MysqlSystem,
			[]sql.Hook{trace.NewTraceHook()},
			&mysql.MySQLDriver{},
			false,
		},
		{
			sql.PostgresSystem,
			[]sql.Hook{log.NewLogHook(), trace.NewTraceHook()},
			&pq.Driver{},
			false,
		},
		{
			sql.UnknownSystem,
			[]sql.Hook{},
			nil,
			true,
		},
		{
			sql.System("invalid"),
			[]sql.Hook{},
			nil,
			true,
		},
	}

	for _, test := range tests {
		d, err := factory.Create(test.system, test.hooks...)

		if test.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.IsType(t, test.expectedBase, d.Base())
			assert.Equal(t, test.hooks, d.Configuration().Hooks())
		}
	}
}
