package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/ankorstore/yokai/fxhttpserver/testdata/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var testMetricsHandlerHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "test_metrics_handler_duration_seconds",
	Help:    "The duration of the TestMetricsHandler",
	Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
})

type TestMetricsHandler struct {
	service *service.TestService
}

func NewTestMetricsHandler(service *service.TestService) *TestMetricsHandler {
	return &TestMetricsHandler{
		service: service,
	}
}

func (h *TestMetricsHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Info("in metrics handler")

		start := time.Now()
		defer func() {
			testMetricsHandlerHistogram.Observe(time.Since(start).Seconds())
		}()

		return c.String(http.StatusOK, fmt.Sprintf("name: %s", h.service.GetAppName()))
	}
}
