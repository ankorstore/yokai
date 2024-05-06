package sql_test

import (
	"context"
	basesql "database/sql"
	"fmt"
	"testing"

	yokailog "github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/log"
	"github.com/ankorstore/yokai/sql/hook/trace"
	yokaitrace "github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func TestRegisterTwice(t *testing.T) {
	driverName1 := registerTestDriver(t)
	driverName2 := registerTestDriver(t)

	expectedDriverName := fmt.Sprintf("%s-sqlite", sql.DriverRegistrationPrefix)
	assert.Equal(t, expectedDriverName, driverName1)
	assert.Equal(t, expectedDriverName, driverName2)
}

func TestRegisterAndPing(t *testing.T) {
	driver := registerTestDriver(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	err = db.Ping()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndPingContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	err = db.PingContext(createTestContext(logger, tracerProvider))
	assert.NoError(t, err)

	logtest.AssertHasNotLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:ping",
		"message":   "sql logger",
	})

	tracetest.AssertHasNotTraceSpan(
		t,
		traceExporter,
		"SQL connection:ping",
		semconv.DBSystemKey.String("sqlite"),
	)

	err = db.Close()
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:close",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:close",
		semconv.DBSystemKey.String("sqlite"),
	)
}

func TestRegisterAndExec(t *testing.T) {
	driver := registerTestDriver(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)")
	assert.NoError(t, err)

	results, err := db.Exec("INSERT INTO foo (bar) VALUES ($1),($2)", "42", "24")
	assert.NoError(t, err)

	lastInsertId, err := results.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), lastInsertId)

	rowsAffected, err := results.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), rowsAffected)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndExecContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	_, err = db.ExecContext(
		createTestContext(logger, tracerProvider),
		"CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)",
	)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "debug",
		"system":       "sqlite",
		"operation":    "connection:exec-context",
		"query":        "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)",
		"lastInsertId": 0,
		"rowsAffected": 0,
		"message":      "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)"),
	)

	results, err := db.ExecContext(
		createTestContext(logger, tracerProvider),
		"INSERT INTO foo (bar) VALUES ($1)",
		42,
	)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "debug",
		"system":       "sqlite",
		"operation":    "connection:exec-context",
		"query":        "INSERT INTO foo (bar) VALUES ($1)",
		"arguments":    "[map[Name: Ordinal:1 Value:42]]",
		"lastInsertId": 1,
		"rowsAffected": 1,
		"message":      "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "INSERT INTO foo (bar) VALUES ($1)"),
		attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:42}]"),
		attribute.Int64("db.lastInsertId", int64(1)),
		attribute.Int64("db.rowsAffected", int64(1)),
	)

	lastInsertId, err := results.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertId)

	rowsAffected, err := results.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndQuery(t *testing.T) {
	driver := registerTestDriver(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	rows, err := db.Query("SELECT $1 AS foo", "bar")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	cols, err := rows.Columns()
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo"}, cols)

	for rows.Next() {
		var foo string
		err = rows.Scan(&foo)
		assert.NoError(t, err)
		assert.Equal(t, "bar", foo)
	}

	err = rows.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndQueryContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	rows, err := db.QueryContext(createTestContext(logger, tracerProvider), "SELECT $1 AS foo", "bar")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:query-context",
		"query":     "SELECT $1 AS foo",
		"arguments": "[map[Name: Ordinal:1 Value:bar]]",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:query-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "SELECT $1 AS foo"),
		attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:bar}]"),
	)

	cols, err := rows.Columns()
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo"}, cols)

	for rows.Next() {
		var foo string
		err = rows.Scan(&foo)
		assert.NoError(t, err)
		assert.Equal(t, "bar", foo)
	}

	err = db.Close()
	assert.NoError(t, err)

	err = rows.Close()
	assert.NoError(t, err)
}

func registerTestDriver(t *testing.T) string {
	t.Helper()

	driver, err := sql.Register(
		"sqlite",
		trace.NewTraceHook(
			trace.WithArguments(true),
			trace.WithExcludedOperations(
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
			),
		),
		log.NewLogHook(
			log.WithArguments(true),
			log.WithExcludedOperations(
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
			),
		),
	)
	assert.NoError(t, err)

	return driver
}

func createTestLogTools(t *testing.T) (*yokailog.Logger, logtest.TestLogBuffer) {
	t.Helper()

	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := yokailog.NewDefaultLoggerFactory().Create(
		yokailog.WithLevel(zerolog.DebugLevel),
		yokailog.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	return logger, logBuffer
}

func createTestTraceTools(t *testing.T) (*otelsdktrace.TracerProvider, tracetest.TestTraceExporter) {
	t.Helper()

	traceExporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := yokaitrace.NewDefaultTracerProviderFactory().Create(
		yokaitrace.WithSpanProcessor(yokaitrace.NewTestSpanProcessor(traceExporter)),
	)
	assert.NoError(t, err)

	return tracerProvider, traceExporter
}

func createTestContext(logger *yokailog.Logger, tracerProvider *otelsdktrace.TracerProvider) context.Context {
	return yokaitrace.WithContext(logger.WithContext(context.Background()), tracerProvider)
}
