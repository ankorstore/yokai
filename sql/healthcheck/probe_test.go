package healthcheck_test

import (
	"context"
	basesql "database/sql"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	t.Parallel()

	db := createTestDB(t)

	probe := healthcheck.NewSQLProbe(db)

	assert.Equal(t, "sql", probe.Name())
}

func TestSetName(t *testing.T) {
	t.Parallel()

	db := createTestDB(t)

	probe := healthcheck.NewSQLProbe(db)

	probe.SetName("custom")

	assert.Equal(t, "custom", probe.Name())
}

func TestCheckSuccess(t *testing.T) {
	t.Parallel()

	db := createTestDB(t)

	probe := healthcheck.NewSQLProbe(db)

	result := probe.Check(context.Background())
	assert.True(t, result.Success)
	assert.Equal(t, "database ping success", result.Message)
}

func TestCheckFailure(t *testing.T) {
	t.Parallel()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.DebugLevel),
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	db := createTestDB(t)

	probe := healthcheck.NewSQLProbe(db)

	err = db.Close()
	assert.NoError(t, err)

	result := probe.Check(logger.WithContext(context.Background()))
	assert.False(t, result.Success)
	assert.Equal(t, "database ping error: sql: database is closed", result.Message)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"error":   "sql: database is closed",
		"message": "database ping error",
	})
}

func createTestDB(t *testing.T) *basesql.DB {
	t.Helper()

	driver, err := sql.Register("sqlite")
	assert.NoError(t, err)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	return db
}
