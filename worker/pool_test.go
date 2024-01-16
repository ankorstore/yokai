package worker_test

import (
	"context"
	"fmt"
	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/ankorstore/yokai/worker"
	"github.com/ankorstore/yokai/worker/testdata/workers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"strings"
	"testing"
	"time"
)

const testExecutionId = "test-execution-id"

func TestWithOneShotWorker(t *testing.T) {
	t.Parallel()

	// test log buffer
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// test trace exporter
	traceExporter := tracetest.NewDefaultTestTraceExporter()
	tracerProvider, err := trace.NewDefaultTracerProviderFactory().Create(
		trace.Global(true),
		trace.WithSpanProcessor(trace.NewTestSpanProcessor(traceExporter)),
	)
	assert.NoError(t, err)

	worker.AnnotateTracerProvider(tracerProvider)

	// test metrics registry
	registry := prometheus.NewPedanticRegistry()

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// pool
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(workers.NewOneShotWorker()),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	err = pool.Metrics().Register(registry)
	assert.NoError(t, err)

	ctx := logger.WithContext(context.Background())
	ctx = context.WithValue(ctx, trace.CtxKey{}, tracerProvider)

	err = pool.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	err = pool.Stop()
	assert.NoError(t, err)

	// post stop execution assertions
	assert.Len(t, pool.Executions(), 1)

	execution, err := pool.Execution("OneShotWorker")
	assert.NoError(t, err)

	assert.Equal(t, worker.Success, execution.Status())
	assert.Equal(t, testExecutionId, execution.Id())
	assert.Len(t, execution.Events(), 2)
	assert.True(t, execution.HasEvent("starting execution attempt 1/1"))
	assert.True(t, execution.HasEvent("stopping execution attempt 1/1 with success"))

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "OneShotWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 1/1",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "OneShotWorker",
		"workerExecutionID": testExecutionId,
		"message":           fmt.Sprintf("running worker OneShotWorker [id %s]", testExecutionId),
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "OneShotWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopped",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "OneShotWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 1/1 with success",
	})

	// traces assertions
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"one shot span",
		attribute.String(worker.TraceSpanAttributeWorkerName, "OneShotWorker"),
		attribute.String(worker.TraceSpanAttributeWorkerExecutionId, testExecutionId),
	)

	// metrics assertions
	expectedMetrics := `
		# HELP worker_execution_total Total number of workers executions
        # TYPE worker_execution_total counter
		worker_execution_total{status="started",workers="oneshotworker"} 1
        worker_execution_total{status="success",workers="oneshotworker"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedMetrics),
		"worker_execution_total",
	)
	assert.NoError(t, err)
}

func TestWithDeferredLoopWorker(t *testing.T) {
	t.Parallel()

	// test log buffer
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// test metrics registry
	registry := prometheus.NewPedanticRegistry()

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// pool
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			workers.NewLoopWorker(),
			worker.WithDeferredStartThreshold(0.1),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	err = pool.Metrics().Register(registry)
	assert.NoError(t, err)

	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	// run execution assertions
	assert.Len(t, pool.Executions(), 1)

	execution, err := pool.Execution("LoopWorker")
	assert.NoError(t, err)

	assert.Equal(t, worker.Running, execution.Status())
	assert.Equal(t, testExecutionId, execution.Id())

	assert.Len(t, execution.Events(), 2)
	assert.True(t, execution.HasEvent("deferring execution attempt for 0.1 seconds"))
	assert.True(t, execution.HasEvent("starting execution attempt 1/1"))

	err = pool.Stop()
	assert.NoError(t, err)

	// post stop execution assertions
	assert.Equal(t, worker.Success, execution.Status())

	assert.Len(t, execution.Events(), 3)
	assert.True(t, execution.HasEvent("stopping execution attempt 1/1 with success"))

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "deferring execution attempt for 0.1 seconds",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 1/1",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopped",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "LoopWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 1/1 with success",
	})

	// metrics assertions
	expectedMetrics := `
		# HELP worker_execution_total Total number of workers executions
        # TYPE worker_execution_total counter
		worker_execution_total{status="started",workers="loopworker"} 1
        worker_execution_total{status="success",workers="loopworker"} 1
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedMetrics),
		"worker_execution_total",
	)
	assert.NoError(t, err)
}

func TestWithRestartingErrorWorker(t *testing.T) {
	t.Parallel()

	// test log buffer
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// test metrics registry
	registry := prometheus.NewPedanticRegistry()

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// pool
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			workers.NewErrorWorker(),
			worker.WithMaxExecutionsAttempts(2),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	err = pool.Metrics().Register(registry)
	assert.NoError(t, err)

	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	time.Sleep(30 * time.Millisecond)

	err = pool.Stop()
	assert.NoError(t, err)

	// post stop execution assertions
	assert.Len(t, pool.Executions(), 1)

	execution, err := pool.Execution("ErrorWorker")
	assert.NoError(t, err)
	assert.Equal(t, worker.Error, execution.Status())

	assert.Len(t, execution.Events(), 6)
	assert.True(t, execution.HasEvent("starting execution attempt 1/2"))
	assert.True(t, execution.HasEvent("stopping execution attempt 1/2 with error: custom error"))
	assert.True(t, execution.HasEvent("restarting after error"))
	assert.True(t, execution.HasEvent("starting execution attempt 2/2"))
	assert.True(t, execution.HasEvent("stopping execution attempt 2/2 with error: custom error"))
	assert.True(t, execution.HasEvent("max execution attempts reached"))

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 1/2",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "error",
		"error":             "custom error",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 1/2 with error: custom error",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "restarting after error",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 2/2",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "error",
		"error":             "custom error",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 2/2 with error: custom error",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "ErrorWorker",
		"workerExecutionID": testExecutionId,
		"message":           "max execution attempts reached",
	})

	// metrics assertions
	expectedMetrics := `
		# HELP worker_execution_total Total number of workers executions
        # TYPE worker_execution_total counter
        worker_execution_total{status="error",workers="errorworker"} 2
        worker_execution_total{status="restarted",workers="errorworker"} 1
        worker_execution_total{status="started",workers="errorworker"} 2
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedMetrics),
		"worker_execution_total",
	)
	assert.NoError(t, err)
}

func TestWithRestartingPanicWorker(t *testing.T) {
	t.Parallel()

	// test log buffer
	logBuffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(logBuffer),
	)
	assert.NoError(t, err)

	// test metrics registry
	registry := prometheus.NewPedanticRegistry()

	// test generator
	generator := uuid.NewTestUuidGenerator(testExecutionId)

	// pool
	pool, err := worker.NewDefaultWorkerPoolFactory().Create(
		worker.WithWorker(
			workers.NewPanicWorker(),
			worker.WithMaxExecutionsAttempts(2),
		),
		worker.WithGenerator(generator),
	)
	assert.NoError(t, err)

	err = pool.Metrics().Register(registry)
	assert.NoError(t, err)

	err = pool.Start(logger.WithContext(context.Background()))
	assert.NoError(t, err)

	time.Sleep(30 * time.Millisecond)

	err = pool.Stop()
	assert.NoError(t, err)

	// post stop execution assertions
	assert.Len(t, pool.Executions(), 1)

	execution, err := pool.Execution("PanicWorker")
	assert.NoError(t, err)
	assert.Equal(t, worker.Error, execution.Status())

	assert.Len(t, execution.Events(), 6)
	assert.True(t, execution.HasEvent("starting execution attempt 1/2"))
	assert.True(t, execution.HasEvent("stopping execution attempt 1/2 with recovered panic: custom panic"))
	assert.True(t, execution.HasEvent("restarting after panic recovery"))
	assert.True(t, execution.HasEvent("starting execution attempt 2/2"))
	assert.True(t, execution.HasEvent("stopping execution attempt 2/2 with recovered panic: custom panic"))
	assert.True(t, execution.HasEvent("max execution attempts reached"))

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 1/2",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "error",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 1/2 with recovered panic: custom panic",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "restarting after panic recovery",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "starting execution attempt 2/2",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "running",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "error",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "stopping execution attempt 1/2 with recovered panic: custom panic",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":             "info",
		"worker":            "PanicWorker",
		"workerExecutionID": testExecutionId,
		"message":           "max execution attempts reached",
	})

	// metrics assertions
	expectedMetrics := `
		# HELP worker_execution_total Total number of workers executions
        # TYPE worker_execution_total counter
        worker_execution_total{status="error",workers="panicworker"} 2
        worker_execution_total{status="restarted",workers="panicworker"} 1
        worker_execution_total{status="started",workers="panicworker"} 2
	`

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expectedMetrics),
		"worker_execution_total",
	)
	assert.NoError(t, err)
}
