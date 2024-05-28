package fxorm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxorm"
	"github.com/ankorstore/yokai/fxorm/testdata/factory"
	"github.com/ankorstore/yokai/fxorm/testdata/model"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

type fakeT struct {
	t      *testing.T
	errors []string
}

func (c *fakeT) Errorf(format string, args ...interface{}) {
	c.errors = append(c.errors, fmt.Sprintf(format, args...))
}

func (c *fakeT) Logf(format string, args ...interface{}) {
	// noop
}

func (c *fakeT) FailNow() {
	// noop
}

func TestModuleWithSqliteAndWithLogEnabledWithValuesAndWithTracesEnabledWithValues(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ORM_DRIVER", "sqlite")
	t.Setenv("ORM_DSN", ":memory:")
	t.Setenv("ORM_LOG_ENABLED", "true")
	t.Setenv("ORM_LOG_LEVEL", "info")
	t.Setenv("ORM_LOG_VALUES", "true")
	t.Setenv("ORM_TRACE_ENABLED", "true")
	t.Setenv("ORM_TRACE_VALUES", "true")

	ctx := context.Background()

	var gormDB *gorm.DB
	var repository *model.TestModelRepository
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,
		fxorm.RunFxOrmAutoMigrate(&model.TestModel{}),
		fx.Provide(model.NewModelRepository),
		fx.Invoke(func(logger *log.Logger, repository *model.TestModelRepository) {
			_ = repository.Create(logger.WithContext(context.Background()), &model.TestModel{
				Name: "test",
			})
		}),
		fx.Populate(&gormDB, &repository, &logBuffer, &traceExporter),
	).RequireStart().RequireStop()

	// assert on DB insertion
	models, err := repository.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, "test", models[0].Name)

	// assert on SQL logs
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "starting ORM auto migration",
	})
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "ORM auto migration success",
	})
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "debug",
		"service":  "test",
		"sqlQuery": "INSERT INTO `test_models` (`name`) VALUES (\"test\")",
	})

	// assert on SQL traces
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"orm.Create",
		semconv.DBSystemKey.String("sqlite"),
		semconv.DBStatementKey.String("INSERT INTO `test_models` (`name`) VALUES (\"test\")"),
	)

	// close
	db, err := gormDB.DB()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestModuleWithSqliteAndWithLogEnabledWithoutValuesAndWithTracesEnabledWithoutValues(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ORM_DRIVER", "sqlite")
	t.Setenv("ORM_DSN", ":memory:")
	t.Setenv("ORM_LOG_ENABLED", "true")
	t.Setenv("ORM_LOG_LEVEL", "info")
	t.Setenv("ORM_LOG_VALUES", "false")
	t.Setenv("ORM_TRACE_ENABLED", "true")
	t.Setenv("ORM_TRACE_VALUES", "false")

	ctx := context.Background()

	var gormDB *gorm.DB
	var repository *model.TestModelRepository
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,
		fxorm.RunFxOrmAutoMigrate(&model.TestModel{}),
		fx.Provide(model.NewModelRepository),
		fx.Invoke(func(logger *log.Logger, repository *model.TestModelRepository) {
			_ = repository.Create(logger.WithContext(context.Background()), &model.TestModel{
				Name: "test",
			})
		}),
		fx.Populate(&gormDB, &repository, &logBuffer, &traceExporter),
	).RequireStart().RequireStop()

	// assert on DB insertion
	models, err := repository.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, "test", models[0].Name)

	// assert on SQL logs
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "starting ORM auto migration",
	})
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "ORM auto migration success",
	})
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":    "debug",
		"service":  "test",
		"sqlQuery": "INSERT INTO `test_models` (`name`) VALUES (?)",
	})

	// assert on SQL traces
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"orm.Create",
		semconv.DBSystemKey.String("sqlite"),
		semconv.DBStatementKey.String("INSERT INTO `test_models` (`name`) VALUES (\"?\")"),
	)

	// close
	db, err := gormDB.DB()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestModuleWithSqliteAndWithLogDisabledAndWithTracesDisabled(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ORM_DRIVER", "sqlite")
	t.Setenv("ORM_DSN", ":memory:")
	t.Setenv("ORM_LOG_ENABLED", "false")
	t.Setenv("ORM_TRACE_ENABLED", "false")

	ctx := context.Background()

	var gormDB *gorm.DB
	var repository *model.TestModelRepository
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,
		fxorm.RunFxOrmAutoMigrate(&model.TestModel{}),
		fx.Provide(model.NewModelRepository),
		fx.Invoke(func(logger *log.Logger, repository *model.TestModelRepository) {
			_ = repository.Create(logger.WithContext(context.Background()), &model.TestModel{
				Name: "test",
			})
		}),
		fx.Populate(&gormDB, &repository, &logBuffer, &traceExporter),
	).RequireStart().RequireStop()

	// assert on DB insertion
	models, err := repository.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, "test", models[0].Name)

	// assert on SQL logs
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "starting ORM auto migration",
	})
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "test",
		"message": "ORM auto migration success",
	})
	hasRecord, _ := logBuffer.HasRecord(map[string]interface{}{
		"level":    "debug",
		"service":  "test",
		"sqlQuery": "INSERT INTO `test_models` (`name`) VALUES (?)",
	})
	assert.False(t, hasRecord)

	// assert on SQL traces
	assert.False(t, traceExporter.HasSpan("orm.Create"))

	// close
	db, err := gormDB.DB()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestModuleWithAutoMigrationError(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("ORM_DRIVER", "sqlite")
	t.Setenv("ORM_DSN", ":memory:")

	ft := &fakeT{t: t}

	fxtest.New(
		ft,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,
		fxorm.RunFxOrmAutoMigrate(&struct {
			ID uint `gorm:"-:all"`
		}{}),
	).RequireStart().RequireStop()

	assert.NotEmpty(t, ft.errors)
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var gormDB *gorm.DB

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fxorm.FxOrmModule,
		fx.Decorate(factory.NewTestOrmFactory),
		fx.Populate(&gormDB),
	).RequireStart().RequireStop()

	assert.True(t, gormDB.DryRun)

	// close
	db, err := gormDB.DB()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
