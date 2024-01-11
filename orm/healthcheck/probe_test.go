package healthcheck_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/orm"
	"github.com/ankorstore/yokai/orm/healthcheck"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	t.Parallel()

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
	)
	assert.NoError(t, err)

	probe := healthcheck.NewOrmProbe(db)

	assert.Equal(t, "orm", probe.Name())
}

func TestSetName(t *testing.T) {
	t.Parallel()

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
	)
	assert.NoError(t, err)

	probe := healthcheck.NewOrmProbe(db)
	probe.SetName("custom")

	assert.Equal(t, "custom", probe.Name())
}

func TestCheckSuccess(t *testing.T) {
	t.Parallel()

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
	)
	assert.NoError(t, err)

	probe := healthcheck.NewOrmProbe(db)

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

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
	)
	assert.NoError(t, err)

	probe := healthcheck.NewOrmProbe(db)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Close()
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
