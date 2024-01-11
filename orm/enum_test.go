package orm_test

import (
	"testing"

	"github.com/ankorstore/yokai/orm"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)

func TestDriverAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		driver   orm.Driver
		expected string
	}{
		{orm.Sqlite, "sqlite"},
		{orm.Mysql, "mysql"},
		{orm.Postgres, "postgres"},
		{orm.SqlServer, "sqlserver"},
		{orm.Unknown, "unknown"},
		{orm.Driver(999), "unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.driver.String())
	}
}

func TestFetchDriver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected orm.Driver
	}{
		{"sqlite", orm.Sqlite},
		{"mysql", orm.Mysql},
		{"postgres", orm.Postgres},
		{"sqlserver", orm.SqlServer},
		{"unknown", orm.Unknown},
		{"", orm.Unknown},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, orm.FetchDriver(tt.input))
	}
}

func TestFetchLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected logger.LogLevel
	}{
		{"silent", logger.Silent},
		{"info", logger.Info},
		{"warn", logger.Warn},
		{"error", logger.Error},
		{"unknown", logger.Silent},
		{"", logger.Silent},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, orm.FetchLogLevel(tt.input))
	}
}
