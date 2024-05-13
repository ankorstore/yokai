package fxsql_test

import (
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxsql/testdata/seed"
	yokailog "github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSeederRunError(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("testdata/config"),
	)
	assert.NoError(t, err)

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := yokailog.NewDefaultLoggerFactory().Create(
		yokailog.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	seeder := fxsql.NewSeeder(db, logger, seed.NewTestSeed(cfg))

	err = db.Close()
	assert.NoError(t, err)

	err = seeder.Run(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "sql: database is closed", err.Error())

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "error",
		"seed":    "test",
		"error":   "sql: database is closed",
		"message": "seed error",
	})
}
