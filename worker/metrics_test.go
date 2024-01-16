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

	metrics := worker.NewWorkerMetrics("foo", "bar")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	metrics.IncrementWorkerExecutionStart("foo")
	metrics.IncrementWorkerExecutionRestart("foo")
	metrics.IncrementWorkerExecutionError("foo")
	metrics.IncrementWorkerExecutionSuccess("foo")

	expected := `
		# HELP foo_bar_worker_execution_total Total number of workers executions
        # TYPE foo_bar_worker_execution_total counter
        foo_bar_worker_execution_total{status="error",worker="foo"} 1
        foo_bar_worker_execution_total{status="restarted",worker="foo"} 1
        foo_bar_worker_execution_total{status="started",worker="foo"} 1
        foo_bar_worker_execution_total{status="success",worker="foo"} 1
	`
	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"foo_bar_worker_execution_total",
	)
	assert.NoError(t, err)
}

func TestWorkerMetricsAlreadyRegistered(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := worker.NewWorkerMetrics("foo", "bar")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	err = metrics.Register(registry)
	assert.Error(t, err)
	assert.Equal(t, "duplicate metrics collector registration attempted", err.Error())
}
