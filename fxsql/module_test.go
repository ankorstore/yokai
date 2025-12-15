package fxsql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxsql/testdata/hook"
	"github.com/ankorstore/yokai/fxsql/testdata/seed"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	var ctx context.Context
	var db *sql.DB
	var logBuffer logtest.TestLogBuffer
	var tracerProvider oteltrace.TracerProvider
	var traceExporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		// provide test hook
		fxsql.AsSQLHook(hook.NewTestHook),
		// provide test seeder
		fxsql.AsSQLSeed(seed.NewTestSeed),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations
		fxsql.RunFxSQLMigration("up"),
		// apply valid seed
		fxsql.RunFxSQLSeeds(),
		// populate test components
		fx.Populate(&ctx, &db, &logBuffer, &tracerProvider, &traceExporter),
	).RequireStart().RequireStop()

	ctx = trace.WithContext(ctx, tracerProvider)

	// SQL query
	rows, err := db.QueryContext(ctx, "SELECT bar FROM foo LIMIT 1")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	// SQL query observability assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:query-context",
		"query":     "SELECT bar FROM foo LIMIT 1",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:query-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "SELECT bar FROM foo LIMIT 1"),
	)

	// SQL query result assertion
	for rows.Next() {
		var bar string
		err = rows.Scan(&bar)
		assert.NoError(t, err)
		assert.Equal(t, "test seed value", bar)
	}

	// SQL seed assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"seed":    "test",
		"message": "seed success",
	})

	// SQL hook assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test hook before connection:exec-context",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test hook after connection:exec-context",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test hook before connection:query-context",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "test hook after connection:query-context",
	})

	// SQL exec
	res, err := db.ExecContext(ctx, "DELETE FROM foo WHERE bar = ?", "test seed value")
	assert.NoError(t, err)

	// SQL exec result assertion
	rowsAffected, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// SQL exec observability assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:exec-context",
		"query":     "DELETE FROM foo WHERE bar = ?",
		"arguments": "[map[Name: Ordinal:1 Value:test seed value]]",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "DELETE FROM foo WHERE bar = ?"),
		attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:test seed value}]"),
	)

	// SQL ping
	err = db.PingContext(ctx)
	assert.NoError(t, err)

	// SQL ping observability assertion (should be excluded)
	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:ping",
		"message":   "sql logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"SQL ping",
		semconv.DBSystemKey.String("sqlite"),
	)

	// SQL close
	err = db.Close()
	assert.NoError(t, err)
}

func TestModuleWithMigrationShutdown(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	ctx := context.Background()

	app := fx.New(
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return ctx
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations and shutdown
		fxsql.RunFxSQLMigrationAndShutdown("up"),
	)

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)
}

func TestModuleErrorWithInvalidDriver(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "invalid")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "invalid")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "invalid")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	ctx := context.Background()

	app := fx.New(
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return ctx
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations
		fxsql.RunFxSQLMigration("up"),
	)

	err := app.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database system for driver invalid")
}

func TestModuleErrorWithInvalidDsn(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "mysql")
	t.Setenv("SQL_PRIMARY_DSN", "invalid")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "mysql")
	t.Setenv("SQL_AUXILIARY1_DSN", "invalid")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "mysql")
	t.Setenv("SQL_AUXILIARY2_DSN", "invalid")

	ctx := context.Background()

	app := fx.New(
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return ctx
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations
		fxsql.RunFxSQLMigration("up"),
	)

	err := app.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid DSN: missing the slash separating the database name")
}

func TestModuleErrorWithInvalidSeed(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	ctx := context.Background()

	var db *sql.DB

	fxtest.New(
		t,
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return ctx
		}),
		// provide test seed
		fxsql.AsSQLSeed(seed.NewTestSeed),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations
		fxsql.RunFxSQLMigration("up"),
		// apply invalid seed
		fxsql.RunFxSQLSeeds("test"),
		// populate test components
		fx.Populate(&db),
	).RequireStart().RequireStop()

	// SQL query
	row := db.QueryRow("SELECT COUNT(*) FROM foo")
	assert.NoError(t, row.Err())

	var count int
	err := row.Scan(&count)
	assert.NoError(t, err)

	// must be 1 (from valid seed, since invalid seed should have been roll backed)
	assert.Equal(t, 1, count)

	// SQL close
	err = db.Close()
	assert.NoError(t, err)
}

func TestModuleWithSeedsShutdown(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	ctx := context.Background()

	app := fx.New(
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return ctx
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply valid seed
		fxsql.RunFxSQLSeedsAndShutdown(),
	)

	err := app.Start(ctx)
	assert.NoError(t, err)

	err = app.Stop(ctx)
	assert.NoError(t, err)
}

func TestModuleDatabasePoolWithDedicatedConnections(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	var pool *fxsql.DatabasePool
	var db *sql.DB

	fxtest.New(
		t,
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// populate test components
		fx.Populate(&pool, &db),
	).RequireStart().RequireStop()

	// verify pool is not nil
	assert.NotNil(t, pool)

	// verify primary database from pool
	assert.NotNil(t, pool.Primary())
	assert.Equal(t, fxsql.PrimaryDatabaseName, pool.Primary().Name())

	// verify primary *sql.DB injection matches pool primary
	assert.Same(t, pool.Primary().DB(), db)

	// verify auxiliaries exist
	aux1, err := pool.Auxiliary("auxiliary1")
	assert.NoError(t, err)
	assert.NotNil(t, aux1)
	assert.Equal(t, "auxiliary1", aux1.Name())

	aux2, err := pool.Auxiliary("auxiliary2")
	assert.NoError(t, err)
	assert.NotNil(t, aux2)
	assert.Equal(t, "auxiliary2", aux2.Name())

	// verify auxiliaries map
	auxiliaries := pool.Auxiliaries()
	assert.Len(t, auxiliaries, 2)
	assert.Contains(t, auxiliaries, "auxiliary1")
	assert.Contains(t, auxiliaries, "auxiliary2")

	// verify each database connection is dedicated by creating unique tables
	// primary database
	_, err = pool.Primary().DB().Exec("CREATE TABLE primary_table (id INTEGER PRIMARY KEY, value TEXT)")
	assert.NoError(t, err)
	_, err = pool.Primary().DB().Exec("INSERT INTO primary_table (value) VALUES ('primary_value')")
	assert.NoError(t, err)

	// auxiliary1 database
	_, err = aux1.DB().Exec("CREATE TABLE aux1_table (id INTEGER PRIMARY KEY, value TEXT)")
	assert.NoError(t, err)
	_, err = aux1.DB().Exec("INSERT INTO aux1_table (value) VALUES ('aux1_value')")
	assert.NoError(t, err)

	// auxiliary2 database
	_, err = aux2.DB().Exec("CREATE TABLE aux2_table (id INTEGER PRIMARY KEY, value TEXT)")
	assert.NoError(t, err)
	_, err = aux2.DB().Exec("INSERT INTO aux2_table (value) VALUES ('aux2_value')")
	assert.NoError(t, err)

	// verify primary only has primary_table
	var primaryValue string
	err = pool.Primary().DB().QueryRow("SELECT value FROM primary_table").Scan(&primaryValue)
	assert.NoError(t, err)
	assert.Equal(t, "primary_value", primaryValue)

	// verify primary does NOT have aux1_table
	err = pool.Primary().DB().QueryRow("SELECT value FROM aux1_table").Scan(&primaryValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such table")

	// verify aux1 only has aux1_table
	var aux1Value string
	err = aux1.DB().QueryRow("SELECT value FROM aux1_table").Scan(&aux1Value)
	assert.NoError(t, err)
	assert.Equal(t, "aux1_value", aux1Value)

	// verify aux1 does NOT have primary_table
	err = aux1.DB().QueryRow("SELECT value FROM primary_table").Scan(&aux1Value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such table")

	// verify aux2 only has aux2_table
	var aux2Value string
	err = aux2.DB().QueryRow("SELECT value FROM aux2_table").Scan(&aux2Value)
	assert.NoError(t, err)
	assert.Equal(t, "aux2_value", aux2Value)

	// verify aux2 does NOT have aux1_table or primary_table
	err = aux2.DB().QueryRow("SELECT value FROM aux1_table").Scan(&aux2Value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such table")

	err = aux2.DB().QueryRow("SELECT value FROM primary_table").Scan(&aux2Value)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such table")

	// cleanup
	err = pool.Primary().DB().Close()
	assert.NoError(t, err)
	err = aux1.DB().Close()
	assert.NoError(t, err)
	err = aux2.DB().Close()
	assert.NoError(t, err)
}

func TestModuleDatabasePoolAuxiliaryNotFound(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_PRIMARY_DRIVER", "sqlite")
	t.Setenv("SQL_PRIMARY_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY1_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY1_DSN", ":memory:")
	t.Setenv("SQL_AUXILIARY2_DRIVER", "sqlite")
	t.Setenv("SQL_AUXILIARY2_DSN", ":memory:")

	var pool *fxsql.DatabasePool

	fxtest.New(
		t,
		fx.NopLogger,
		// provide context
		fx.Provide(func() context.Context {
			return context.Background()
		}),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// populate test components
		fx.Populate(&pool),
	).RequireStart().RequireStop()

	// verify auxiliary not found error
	_, err := pool.Auxiliary("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "database with name nonexistent was not found", err.Error())

	// cleanup
	err = pool.Primary().DB().Close()
	assert.NoError(t, err)
}
