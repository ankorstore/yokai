package fxsql_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMigratorRunErrorWithInvalidDialect(t *testing.T) {
	t.Parallel()

	migrator := createTestMigrator(t)

	err := migrator.Run(context.Background(), "invalid", "", "")
	assert.Error(t, err)
	assert.Equal(t, `"invalid": unknown dialect`, err.Error())
}

func TestMigratorRunErrorWithInvalidCommand(t *testing.T) {
	t.Parallel()

	migrator := createTestMigrator(t)

	err := migrator.Run(context.Background(), "sqlite", "", "invalid")
	assert.Error(t, err)
	assert.Equal(t, `"invalid": no such command`, err.Error())
}

func createTestMigrator(t *testing.T) *fxsql.Migrator {
	t.Helper()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.Disabled),
	)
	assert.NoError(t, err)

	return fxsql.NewMigrator(nil, logger)
}
