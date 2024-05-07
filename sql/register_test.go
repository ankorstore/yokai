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

func TestRegisterAndExecContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	// create table
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

	// insert into table
	result, err := db.ExecContext(
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

	lastInsertId, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertId)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

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

	err = rows.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndPrepareContextAndExecContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	ctx := createTestContext(logger, tracerProvider)

	// create table
	_, err = db.ExecContext(ctx, "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)")
	assert.NoError(t, err)

	// prepare insert into table
	stmt, err := db.PrepareContext(ctx, "INSERT INTO foo (bar) VALUES ($1)")
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:prepare-context",
		"query":     "INSERT INTO foo (bar) VALUES ($1)",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:prepare-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "INSERT INTO foo (bar) VALUES ($1)"),
	)

	// exec insert into table
	result, err := stmt.ExecContext(ctx, 42)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "debug",
		"system":       "sqlite",
		"operation":    "statement:exec-context",
		"query":        "INSERT INTO foo (bar) VALUES ($1)",
		"arguments":    "[map[Name: Ordinal:1 Value:42]]",
		"lastInsertId": 1,
		"rowsAffected": 1,
		"message":      "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL statement:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "INSERT INTO foo (bar) VALUES ($1)"),
		attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:42}]"),
		attribute.Int64("db.lastInsertId", 1),
		attribute.Int64("db.rowsAffected", 1),
	)

	lastInsertId, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertId)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	err = stmt.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndPrepareContextAndQueryContext(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	ctx := createTestContext(logger, tracerProvider)

	stmt, err := db.PrepareContext(ctx, "SELECT $1 AS foo")
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:prepare-context",
		"query":     "SELECT $1 AS foo",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:prepare-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "SELECT $1 AS foo"),
	)

	rows, err := stmt.QueryContext(ctx, "bar")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "statement:query-context",
		"query":     "SELECT $1 AS foo",
		"arguments": "[map[Name: Ordinal:1 Value:bar]]",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL statement:query-context",
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

	err = rows.Close()
	assert.NoError(t, err)

	err = stmt.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndBeginTxAndCommit(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	ctx := createTestContext(logger, tracerProvider)

	// create table
	_, err = db.ExecContext(ctx, "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)")
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:exec-context",
		"query":     "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)"),
	)

	// begin transaction
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:begin-tx",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:begin-tx",
		semconv.DBSystemKey.String("sqlite"),
	)

	// insert into table
	result, err := tx.ExecContext(ctx, "INSERT INTO foo (bar) VALUES ($1)", 42)
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
		attribute.Int64("db.lastInsertId", 1),
		attribute.Int64("db.rowsAffected", 1),
	)

	lastInsertId, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertId)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// commit transaction
	err = tx.Commit()
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "transaction:commit",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL transaction:commit",
		semconv.DBSystemKey.String("sqlite"),
	)

	// check inserted row is present
	rows, err := db.Query("SELECT count(*) FROM foo")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	}

	err = rows.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestRegisterAndBeginTxAndRollback(t *testing.T) {
	driver := registerTestDriver(t)
	logger, logBuffer := createTestLogTools(t)
	tracerProvider, traceExporter := createTestTraceTools(t)

	db, err := basesql.Open(driver, ":memory:")
	assert.NoError(t, err)

	ctx := createTestContext(logger, tracerProvider)

	// create table
	_, err = db.ExecContext(ctx, "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)")
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:exec-context",
		"query":     "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:exec-context",
		semconv.DBSystemKey.String("sqlite"),
		attribute.String("db.statement", "CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, bar INTEGER)"),
	)

	// begin transaction
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "connection:begin-tx",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL connection:begin-tx",
		semconv.DBSystemKey.String("sqlite"),
	)

	// insert into table
	result, err := tx.ExecContext(ctx, "INSERT INTO foo (bar) VALUES ($1)", 42)
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
		attribute.Int64("db.lastInsertId", 1),
		attribute.Int64("db.rowsAffected", 1),
	)

	lastInsertId, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertId)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// rollback transaction
	err = tx.Rollback()
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":     "debug",
		"system":    "sqlite",
		"operation": "transaction:rollback",
		"message":   "sql logger",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"SQL transaction:rollback",
		semconv.DBSystemKey.String("sqlite"),
	)

	// check inserted row is not present
	rows, err := db.Query("SELECT count(*) FROM foo")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	}

	err = rows.Close()
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
}

func TestRegisterTwice(t *testing.T) {
	driverName1 := registerTestDriver(t)
	driverName2 := registerTestDriver(t)

	expectedDriverName := fmt.Sprintf("%s-sqlite", sql.DriverRegistrationPrefix)
	assert.Equal(t, expectedDriverName, driverName1)
	assert.Equal(t, expectedDriverName, driverName2)
}

func TestRegisterWithUnsupportedDriver(t *testing.T) {
	driver, err := sql.Register("invalid")

	assert.Equal(t, "", driver)
	assert.Error(t, err)
	assert.Equal(t, "unsupported database system for driver invalid", err.Error())
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
