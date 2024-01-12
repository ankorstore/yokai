package fxmetrics_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxmetrics/testdata/metrics"
	"github.com/stretchr/testify/assert"
)

func TestAsMetricsCollector(t *testing.T) {
	t.Parallel()

	result := fxmetrics.AsMetricsCollector(metrics.FxMetricsTestCounter)

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}

func TestAsMetricsCollectors(t *testing.T) {
	t.Parallel()

	result := fxmetrics.AsMetricsCollectors(metrics.FxMetricsTestCounter, metrics.FxMetricsDuplicatedTestCounter)

	assert.Equal(t, "fx.optionGroup", fmt.Sprintf("%T", result))
}
