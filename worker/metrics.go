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

type WorkerMetrics struct {
	registered bool
	namespace  string
	subsystem  string
	counter    *prometheus.CounterVec
}

func NewWorkerMetrics(namespace string, subsystem string) *WorkerMetrics {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Sanitize(namespace),
			Subsystem: Sanitize(subsystem),
			Name:      "worker_execution_total",
			Help:      "Total number of workers executions",
		},
		[]string{
			"workers",
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

func (m *WorkerMetrics) Register(registry *prometheus.Registry) error {
	err := registry.Register(m.counter)
	if err != nil {
		return err
	}

	m.registered = err == nil

	return err
}

func (m *WorkerMetrics) IncrementWorkerExecutionStart(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionStarted).Inc()
	}

	return m
}

func (m *WorkerMetrics) IncrementWorkerExecutionRestart(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionRestarted).Inc()
	}

	return m
}

func (m *WorkerMetrics) IncrementWorkerExecutionSuccess(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionSuccess).Inc()
	}

	return m
}

func (m *WorkerMetrics) IncrementWorkerExecutionError(workerName string) *WorkerMetrics {
	if m.registered {
		m.counter.WithLabelValues(Sanitize(workerName), ExecutionError).Inc()
	}

	return m
}
