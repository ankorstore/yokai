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

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.system.String())
	}
}

func TestFetchSystemFromDriver(t *testing.T) {
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
			"mysql",
			sql.MysqlSystem,
		},
		{
			"postgres",
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

	for _, tt := range tests {
		assert.Equal(t, tt.expected, sql.FetchSystem(tt.driver))
	}
}
