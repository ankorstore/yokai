package sql_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
)

func TestSystemAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		system   sql.System
		expected string
	}{
		{
			sql.SqliteSystem,
			"sqlite",
		},
		{
			sql.MysqlSystem,
			"mysql",
		},
		{
			sql.PostgresSystem,
			"postgres",
		},
		{
			sql.UnknownSystem,
			"unknown",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.system.String())
	}
}

func TestFetchSystem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		driver   string
		expected sql.System
	}{
		{
			"sqlite",
			sql.SqliteSystem,
		},
		{
			"SQLite",
			sql.SqliteSystem,
		},
		{
			"mysql",
			sql.MysqlSystem,
		},
		{
			"MySQL",
			sql.MysqlSystem,
		},
		{
			"postgres",
			sql.PostgresSystem,
		},
		{
			"Postgres",
			sql.PostgresSystem,
		},
		{
			"",
			sql.UnknownSystem,
		},
		{
			"unknown",
			sql.UnknownSystem,
		},
		{
			"invalid",
			sql.UnknownSystem,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, sql.FetchSystem(test.driver))
	}
}
