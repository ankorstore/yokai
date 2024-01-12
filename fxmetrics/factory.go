package fxmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsRegistryFactory is the interface for [prometheus.Registry] factories.
type MetricsRegistryFactory interface {
	Create() (*prometheus.Registry, error)
}

// DefaultMetricsRegistryFactory is the default [MetricsRegistryFactory] implementation.
type DefaultMetricsRegistryFactory struct{}

// NewDefaultMetricsRegistryFactory returns a [DefaultMetricsRegistryFactory], implementing [MetricsRegistryFactory].
func NewDefaultMetricsRegistryFactory() MetricsRegistryFactory {
	return &DefaultMetricsRegistryFactory{}
}

// Create returns a new [prometheus.Registry].
func (f *DefaultMetricsRegistryFactory) Create() (*prometheus.Registry, error) {
	return prometheus.NewRegistry(), nil
}
