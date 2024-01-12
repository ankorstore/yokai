package fxmetrics_test

import (
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxmetrics/testdata/factory"
	"github.com/ankorstore/yokai/fxmetrics/testdata/metrics"
	"github.com/ankorstore/yokai/fxmetrics/testdata/spy"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var logBuffer logtest.TestLogBuffer
	var registry *prometheus.Registry

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,
		fx.Options(
			fxmetrics.AsMetricsCollector(metrics.FxMetricsTestCounter),
		),
		fx.Invoke(func() {
			metrics.FxMetricsTestCounter.Add(9)
		}),
		fx.Populate(&logBuffer, &registry),
	).RequireStart().RequireStop()

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "debug",
		"message": "registered metrics collector *prometheus.counter",
	})

	expectedHelp := `
		# HELP test_total test help
		# TYPE test_total counter
	`
	expectedMetric := `
		test_total 9
	`

	err := testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedHelp+expectedMetric),
		"test_total",
	)
	assert.NoError(t, err)
}

func TestModuleErrorWithDuplicatedCollector(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var logBuffer logtest.TestLogBuffer
	var registry *prometheus.Registry

	spyTB := spy.NewSpyTB()

	fxtest.New(
		spyTB,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,
		fx.Options(
			fxmetrics.AsMetricsCollectors(metrics.FxMetricsTestCounter, metrics.FxMetricsDuplicatedTestCounter),
		),
		fx.Populate(&logBuffer, &registry),
	).RequireStart().RequireStop()

	assert.Empty(t, spyTB.Logs())

	assert.NotZero(t, spyTB.Failures())
	assert.Contains(t, spyTB.Errors().String(), "duplicate metrics collector registration attempted")
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var logBuffer logtest.TestLogBuffer
	var registry *prometheus.Registry

	spyTB := spy.NewSpyTB()

	fxtest.New(
		spyTB,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxmetrics.FxMetricsModule,
		fx.Decorate(factory.NewTestMetricsRegistryFactory),
		fx.Populate(&logBuffer, &registry),
	).RequireStart().RequireStop()

	assert.Contains(t, spyTB.Logs().String(), "NewTestMetricsRegistryFactory")

	assert.NotZero(t, spyTB.Failures())
	assert.Contains(t, spyTB.Errors().String(), "custom error")
}
