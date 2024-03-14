package worker_test

import (
	"strings"
	"testing"

	"github.com/ankorstore/yokai/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWorkerMetrics(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := worker.NewWorkerMetrics("", "")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	metrics.IncrementWorkerExecutionStart("foo")
	metrics.IncrementWorkerExecutionRestart("foo")
	metrics.IncrementWorkerExecutionError("foo")
	metrics.IncrementWorkerExecutionSuccess("foo")

	expected := `
		# HELP worker_executions_total Total number of workers executions
        # TYPE worker_executions_total counter
        worker_executions_total{status="error",worker="foo"} 1
        worker_executions_total{status="restarted",worker="foo"} 1
        worker_executions_total{status="started",worker="foo"} 1
        worker_executions_total{status="success",worker="foo"} 1
	`
	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"worker_executions_total",
	)
	assert.NoError(t, err)
}

func TestWorkerMetricsWithNamespaceAndSubsystem(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	metrics.IncrementWorkerExecutionStart("foo")
	metrics.IncrementWorkerExecutionRestart("foo")
	metrics.IncrementWorkerExecutionError("foo")
	metrics.IncrementWorkerExecutionSuccess("foo")

	expected := `
		# HELP foo_bar_worker_executions_total Total number of workers executions
        # TYPE foo_bar_worker_executions_total counter
        foo_bar_worker_executions_total{status="error",worker="foo"} 1
        foo_bar_worker_executions_total{status="restarted",worker="foo"} 1
        foo_bar_worker_executions_total{status="started",worker="foo"} 1
        foo_bar_worker_executions_total{status="success",worker="foo"} 1
	`
	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"foo_bar_worker_executions_total",
	)
	assert.NoError(t, err)
}

func TestWorkerMetricsWithCollectorAlreadyRegistered(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	err = metrics.Register(registry)
	assert.Error(t, err)
	assert.Equal(t, "duplicate metrics collector registration attempted", err.Error())
}
