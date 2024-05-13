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
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_DRIVER", "sqlite")
	t.Setenv("SQL_DSN", ":memory:")

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
	t.Setenv("SQL_DRIVER", "sqlite")
	t.Setenv("SQL_DSN", ":memory:")

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
}

func TestModuleErrorWithInvalidDriver(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_DRIVER", "invalid")

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
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_DRIVER", "mysql")
	t.Setenv("SQL_DSN", "invalid")

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
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("SQL_DRIVER", "sqlite")
	t.Setenv("SQL_DSN", ":memory:")

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
