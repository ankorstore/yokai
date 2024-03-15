package fxworker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxgenerate"
	"github.com/ankorstore/yokai/fxgenerate/fxgeneratetest/uuid"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/fxworker"
	"github.com/ankorstore/yokai/fxworker/testdata/factory"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// workerExecutionId is a deterministic test worker execution id.
const workerExecutionId = "testWorkerExecutionID"

func TestModuleWithDefaultConfig(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxworker.FxWorkerModule,
		// deterministic test worker execution id
		fx.Provide(
			fx.Annotate(
				func() string {
					return workerExecutionId
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		fx.Options(
			// workers registration
			fxworker.AsWorker(workers.NewClassicWorker),
		),
		// deterministic generator for worker execution id
		fx.Decorate(uuid.NewFxTestUuidGeneratorFactory),
		// extraction
		fx.Populate(&logBuffer, &traceExporter, &metricsRegistry),
		// invoke worker pool
		fx.Invoke(func(*worker.WorkerPool) {}),
	).RequireStart()

	// 1 seconds for workers to run
	time.Sleep(1 * time.Second)

	app.RequireStop()

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "ClassicWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "starting execution attempt 1/1",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "ClassicWorker",
		"workerExecutionID": workerExecutionId,
		"message":           fmt.Sprintf("running worker ClassicWorker [id %s]", workerExecutionId),
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "ClassicWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "stopping execution attempt 1/1 with success",
	})

	// workers metrics assertions
	expected := `
		# HELP worker_executions_total Total number of workers executions
		# TYPE worker_executions_total counter
        worker_executions_total{status="started",worker="classicworker"} 1
        worker_executions_total{status="success",worker="classicworker"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expected),
		"worker_executions_total",
	)
	assert.NoError(t, err)
}

func TestModuleWithCustomConfig(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "custom")

	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxworker.FxWorkerModule,
		// deterministic test worker execution id
		fx.Provide(
			fx.Annotate(
				func() string {
					return workerExecutionId
				},
				fx.ResultTags(`name:"generate-test-uuid-value"`),
			),
		),
		fx.Options(
			// workers registration
			fxworker.AsWorker(workers.NewCancellableWorker),
		),
		// deterministic generator for worker execution id
		fx.Decorate(uuid.NewFxTestUuidGeneratorFactory),
		// extraction
		fx.Populate(&logBuffer, &traceExporter, &metricsRegistry),
		// invoke worker pool
		fx.Invoke(func(*worker.WorkerPool) {}),
	).RequireStart()

	// 1 seconds for workers to run
	time.Sleep(1 * time.Second)

	app.RequireStop()

	// worker LoopWorker assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "CancellableWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "deferring execution attempt for 0.1 seconds",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "CancellableWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "CancellableWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "stopping",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"service":           "test",
		"module":            "worker",
		"worker":            "CancellableWorker",
		"workerExecutionID": workerExecutionId,
		"message":           "stopped",
	})

	// workers metrics assertions
	expected := `
		# HELP foo_bar_worker_executions_total Total number of workers executions
		# TYPE foo_bar_worker_executions_total counter
        foo_bar_worker_executions_total{status="started",worker="cancellableworker"} 1
        foo_bar_worker_executions_total{status="success",worker="cancellableworker"} 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expected),
		"foo_bar_worker_executions_total",
	)
	assert.NoError(t, err)
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	var pool *worker.WorkerPool

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fxmetrics.FxMetricsModule,
		fxgenerate.FxGenerateModule,
		fxworker.FxWorkerModule,
		fx.Decorate(factory.NewTestWorkerPoolFactory),
		fx.Populate(&pool),
	).RequireStart().RequireStop()

	assert.Equal(t, 99, pool.Options().GlobalMaxExecutionsAttempts)
}
