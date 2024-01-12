package fxmetrics_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMetricsRegistryFactory(t *testing.T) {
	t.Parallel()

	factory := fxmetrics.NewDefaultMetricsRegistryFactory()

	assert.IsType(t, &fxmetrics.DefaultMetricsRegistryFactory{}, factory)
	assert.Implements(t, (*fxmetrics.MetricsRegistryFactory)(nil), factory)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	factory := fxmetrics.NewDefaultMetricsRegistryFactory()

	registry, err := factory.Create()
	assert.NoError(t, err)

	assert.IsType(t, &prometheus.Registry{}, registry)
}
