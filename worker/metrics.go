package worker

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ExecutionStarted   = "started"
	ExecutionRestarted = "restarted"
	ExecutionSuccess   = "success"
	ExecutionError     = "error"
)

// WorkerMetrics allows the [WorkerPool] to send worker metrics to a [prometheus.Registry].
type WorkerMetrics struct {
	registered bool
	namespace  string
	subsystem  string
	counter    *prometheus.CounterVec
}

// NewWorkerMetrics returns a new [WorkerMetrics], and accepts metrics namespace and subsystem.
func NewWorkerMetrics(namespace string, subsystem string) *WorkerMetrics {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Sanitize(namespace),
			Subsystem: Sanitize(subsystem),
			Name:      "worker_execution_total",
			Help:      "Total number of workers executions",
		},
		[]string{
			"worker",
			"status",
		},
	)

	return &WorkerMetrics{
		registered: false,
		namespace:  namespace,
		subsystem:  subsystem,
		counter:    counter,
	}
}

// Register registers the [WorkerMetrics] against a [prometheus.Registry].
func (m *WorkerMetrics) Register(registry *prometheus.Registry) error {
	err := registry.Register(m.counter)
	if err != nil {
		return err
	}

	m.registered = err == nil

	return err
}

// IncrementWorkerExecutionStart increments the started workers counter for a given worker name.
func (m *WorkerMetrics) IncrementWorkerExecutionStart(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionStarted).Inc()
	}

	return m
}

// IncrementWorkerExecutionRestart increments the restarted workers counter for a given worker name.
func (m *WorkerMetrics) IncrementWorkerExecutionRestart(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionRestarted).Inc()
	}

	return m
}

// IncrementWorkerExecutionSuccess increments the successful workers counter for a given worker name.
func (m *WorkerMetrics) IncrementWorkerExecutionSuccess(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionSuccess).Inc()
	}

	return m
}

// IncrementWorkerExecutionError increments the failing workers counter for a given worker name.
func (m *WorkerMetrics) IncrementWorkerExecutionError(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionError).Inc()
	}

	return m
}
