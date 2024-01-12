package factory

import (
	"fmt"

	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
)

type TestMetricsRegistryFactory struct{}

func NewTestMetricsRegistryFactory() fxmetrics.MetricsRegistryFactory {
	return &TestMetricsRegistryFactory{}
}

func (f *TestMetricsRegistryFactory) Create() (*prometheus.Registry, error) {
	return nil, fmt.Errorf("custom error")
}
