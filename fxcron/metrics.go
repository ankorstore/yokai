package fxcron

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	EXECUTION_SUCCESS = "success"
	EXECUTION_ERROR   = "error"
)

var defaultBuckets = []float64{
	// 1ms
	0.001,
	0.002,
	0.005,
	// 10ms
	0.01,
	0.02,
	0.05,
	// 100ms
	0.1,
	0.2,
	0.5,
	// 1s
	1.0,
	2.0,
	5.0,
	// 10s
	10.0,
	20.0,
	50.0,
	// 100s
	100.0,
	200.0,
	500.0,
	// 1000s
	1000.0,
	2000.0,
	5000.0,
}

// CronJobMetrics is the metrics handler for the cron jobs.
type CronJobMetrics struct {
	registered bool
	namespace  string
	subsystem  string
	histogram  *prometheus.HistogramVec
	counter    *prometheus.CounterVec
}

// NewCronJobMetrics returns a new [CronJobMetrics] instance for provided metrics namespace and subsystem.
func NewCronJobMetrics(namespace string, subsystem string) *CronJobMetrics {
	return create(namespace, subsystem, defaultBuckets)
}

// NewCronJobMetricsWithBuckets returns a new [CronJobMetrics] instance for provided metrics namespace, subsystem and buckets.
func NewCronJobMetricsWithBuckets(namespace string, subsystem string, buckets []float64) *CronJobMetrics {
	return create(namespace, subsystem, buckets)
}

// Register allows the [CronJobMetrics] to register against a provided [prometheus.Registry].
func (m *CronJobMetrics) Register(registry *prometheus.Registry) error {
	err := registry.Register(m.histogram)
	if err != nil {
		return err
	}

	err = registry.Register(m.counter)
	if err != nil {
		return err
	}

	m.registered = err == nil

	return err
}

// ObserveCronJobExecutionDuration observes the duration of a cron job execution.
func (m *CronJobMetrics) ObserveCronJobExecutionDuration(jobName string, jobDuration float64) *CronJobMetrics {
	if m.registered {
		m.histogram.WithLabelValues(Sanitize(jobName)).Observe(jobDuration)
	}

	return m
}

// IncrementCronJobExecutionSuccess increments the number of execution successes for a given cron job.
func (m *CronJobMetrics) IncrementCronJobExecutionSuccess(jobName string) *CronJobMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(jobName), EXECUTION_SUCCESS).Inc()
	}

	return m
}

// IncrementCronJobExecutionError increments the number of execution errors for a given cron job.
func (m *CronJobMetrics) IncrementCronJobExecutionError(jobName string) *CronJobMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(jobName), EXECUTION_ERROR).Inc()
	}

	return m
}

func create(namespace string, subsystem string, buckets []float64) *CronJobMetrics {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Sanitize(namespace),
			Subsystem: Sanitize(subsystem),
			Name:      "job_execution_duration_seconds",
			Help:      "Duration of cron job executions in seconds",
			Buckets:   buckets,
		},
		[]string{
			"job",
		},
	)

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Sanitize(namespace),
			Subsystem: Sanitize(subsystem),
			Name:      "job_execution_total",
			Help:      "Total number of cron job executions",
		},
		[]string{
			"job",
			"status",
		},
	)

	return &CronJobMetrics{
		registered: false,
		namespace:  namespace,
		subsystem:  subsystem,
		histogram:  histogram,
		counter:    counter,
	}
}
