package fxcron_test

import (
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxcron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCronJobMetrics(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := fxcron.NewCronJobMetrics("", "")

	err := metrics.Register(registry)
	assert.NoError(t, err)

	// execution duration
	expected := `
		# HELP cron_executions_duration_seconds Duration of cron job executions in seconds
		# TYPE cron_executions_duration_seconds histogram
		cron_executions_duration_seconds_bucket{job="foo",le="0.001"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.002"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.005"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.01"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.02"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.05"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.1"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.2"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="0.5"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="1"} 0
		cron_executions_duration_seconds_bucket{job="foo",le="2"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="5"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="10"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="20"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="50"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="100"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="200"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="500"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="1000"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="2000"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="5000"} 1
		cron_executions_duration_seconds_bucket{job="foo",le="+Inf"} 1
		cron_executions_duration_seconds_sum{job="foo"} 1.1
		cron_executions_duration_seconds_count{job="foo"} 1
	`

	metrics.ObserveCronJobExecutionDuration("foo", 1.1)

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"cron_executions_duration_seconds",
	)
	assert.NoError(t, err)

	// execution counter
	expected = `
		# HELP cron_executions_total Total number of cron job executions
		# TYPE cron_executions_total counter
		cron_executions_total{job="foo",status="success"} 1
		cron_executions_total{job="foo",status="error"} 1
	`

	metrics.IncrementCronJobExecutionSuccess("foo")
	metrics.IncrementCronJobExecutionError("foo")

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"cron_executions_total",
	)
	assert.NoError(t, err)
}

func TestCronJobMetricsWithNamespaceAndSubsystemAbdBuckets(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewPedanticRegistry()

	metrics := fxcron.NewCronJobMetricsWithBuckets("foo", "bar", []float64{0.1, 1, 10})

	err := metrics.Register(registry)
	assert.NoError(t, err)

	// execution duration
	expected := `
		# HELP foo_bar_cron_executions_duration_seconds Duration of cron job executions in seconds
		# TYPE foo_bar_cron_executions_duration_seconds histogram
		foo_bar_cron_executions_duration_seconds_bucket{job="foo",le="0.1"} 0
		foo_bar_cron_executions_duration_seconds_bucket{job="foo",le="1"} 0
		foo_bar_cron_executions_duration_seconds_bucket{job="foo",le="10"} 1
		foo_bar_cron_executions_duration_seconds_bucket{job="foo",le="+Inf"} 1
		foo_bar_cron_executions_duration_seconds_sum{job="foo"} 1.1
		foo_bar_cron_executions_duration_seconds_count{job="foo"} 1
	`

	metrics.ObserveCronJobExecutionDuration("foo", 1.1)

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"foo_bar_cron_executions_duration_seconds",
	)
	assert.NoError(t, err)

	// execution counter
	expected = `
		# HELP foo_bar_cron_executions_total Total number of cron job executions
		# TYPE foo_bar_cron_executions_total counter
		foo_bar_cron_executions_total{job="foo",status="success"} 1
		foo_bar_cron_executions_total{job="foo",status="error"} 1
	`

	metrics.IncrementCronJobExecutionSuccess("foo")
	metrics.IncrementCronJobExecutionError("foo")

	err = testutil.GatherAndCompare(
		registry,
		strings.NewReader(expected),
		"foo_bar_cron_executions_total",
	)
	assert.NoError(t, err)
}
