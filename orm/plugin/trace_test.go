package plugin_test

import (
	"testing"

	"github.com/ankorstore/yokai/orm"
	"github.com/ankorstore/yokai/orm/plugin"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type model struct {
	Name string
}

func TestWithValues(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.NoError(t, err)

	err = db.Use(plugin.NewOrmTracerPlugin(tracerProvider, true))
	assert.NoError(t, err)

	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	err = db.Create(&model{Name: "test"}).Error
	assert.NoError(t, err)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"orm.Create",
		semconv.DBSystemKey.String("sqlite"),
		semconv.DBStatementKey.String("INSERT INTO `models` (`name`) VALUES (\"test\")"),
	)
}

func TestWithoutValues(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.NoError(t, err)

	err = db.Use(plugin.NewOrmTracerPlugin(tracerProvider, false))
	assert.NoError(t, err)

	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	err = db.Create(&model{Name: "test"}).Error
	assert.NoError(t, err)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"orm.Create",
		semconv.DBSystemKey.String("sqlite"),
		semconv.DBStatementKey.String("INSERT INTO `models` (`name`) VALUES (\"?\")"),
	)
}

func TestError(t *testing.T) {
	t.Parallel()

	exporter := tracetest.NewDefaultTestTraceExporter()

	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(exporter)),
	)
	assert.NoError(t, err)

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.NoError(t, err)

	err = db.Use(plugin.NewOrmTracerPlugin(tracerProvider, true))
	assert.NoError(t, err)

	err = db.AutoMigrate(&model{})
	assert.NoError(t, err)

	err = db.Raw("invalid query").Scan(&model{}).Error
	assert.Error(t, err)

	tracetest.AssertHasTraceSpan(
		t,
		exporter,
		"orm.Row",
		semconv.DBSystemKey.String("sqlite"),
		semconv.DBStatementKey.String("invalid query"),
	)
}
