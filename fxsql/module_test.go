package fxsql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/fxsql/testdata/hook"
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

func TestModuleWithObservabilityAndCustomHook(t *testing.T) {
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
		// provide custom hook
		fxsql.AsSQLHook(hook.NewDummyHook),
		// load module and dependencies
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxsql.FxSQLModule,
		// apply migrations
		fxsql.RunFxSQLMigration("up"),
		// populate test components
		fx.Populate(&ctx, &db, &logBuffer, &tracerProvider, &traceExporter),
	).RequireStart().RequireStop()

	ctx = trace.WithContext(ctx, tracerProvider)

	// SQL exec query
	_, err := db.ExecContext(ctx, "INSERT INTO foo (bar) VALUES (?)", "test")
	assert.NoError(t, err)

	// SQL exec observability assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:exec-context",
		"query":     "INSERT INTO foo (bar) VALUES (?)",
		"arguments": "[map[Name: Ordinal:1 Value:test]]",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "INSERT INTO foo (bar) VALUES (?)"),
		attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:test}]"),
	)

	// SQL custom hook assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "DummyHook: before connection:exec-context",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "DummyHook: after connection:exec-context",
	})

	// SQL ping
	err = db.PingContext(ctx)
	assert.NoError(t, err)

	// SQL ping observability assertion (excluded)
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
		// apply migrations and shutdown
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
		// apply migrations and shutdown
		fxsql.RunFxSQLMigration("up"),
	)

	err := app.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid DSN: missing the slash separating the database name")
}
