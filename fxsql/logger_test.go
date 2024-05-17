package fxsql_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMigratorLoggerPrintf(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()

	migratorLogger := createTestMigratorLogger(t, logBuffer)

	migratorLogger.Printf("test %s %d", "foo", 42)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test foo 42",
	})
}

func TestMigratorLoggerFatalf(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()

	migratorLogger := createTestMigratorLogger(t, logBuffer)

	migratorLogger.Fatalf("test %s %d", "foo", 42)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "fatal",
		"message": "test foo 42",
	})
}

func createTestMigratorLogger(t *testing.T, buffer logtest.TestLogBuffer) *fxsql.MigratorLogger {
	t.Helper()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	return fxsql.NewMigratorLogger(logger, true)
}
