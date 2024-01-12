package metrics

import "github.com/prometheus/client_golang/prometheus"

var FxMetricsTestCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "test_total",
	Help: "test help",
})

var FxMetricsDuplicatedTestCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "test_total",
	Help: "test help",
})
