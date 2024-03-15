package fxhttpclient

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/ankorstore/yokai/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	ModuleName                       = "httpclient"
	DefaultTimeout                   = 30
	DefaultMaxIdleConnections        = 100
	DefaultMaxConnectionsPerHost     = 100
	DefaultMaxIdleConnectionsPerHost = 100
)

// FxHttpClientModule is the [Fx] http client module.
//
// [Fx]: https://github.com/uber-go/fx
var FxHttpClientModule = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			NewFxHttpClientTransport,
			fx.As(new(http.RoundTripper)),
		),
		httpclient.NewDefaultHttpClientFactory,
		NewFxHttpClient,
	),
)

// FxHttpClientTransportParam allows injection of the required dependencies in [NewFxHttpClientTransport].
type FxHttpClientTransportParam struct {
	fx.In
	TracerProvider  trace.TracerProvider
	Config          *config.Config
	Logger          *log.Logger
	MetricsRegistry *prometheus.Registry
}

// NewFxHttpClientTransport returns a new [http.RoundTripper].
func NewFxHttpClientTransport(p FxHttpClientTransportParam) http.RoundTripper {
	// base round tripper config
	maxIdleConnections := p.Config.GetInt("modules.http.client.transport.max_idle_connections")
	if maxIdleConnections == 0 {
		maxIdleConnections = DefaultMaxIdleConnections
	}

	maxConnectionsPerHost := p.Config.GetInt("modules.http.client.transport.max_connections_per_host")
	if maxConnectionsPerHost == 0 {
		maxConnectionsPerHost = DefaultMaxConnectionsPerHost
	}

	maxIdleConnectionsPerHost := p.Config.GetInt("modules.http.client.transport.max_idle_connections_per_host")
	if maxIdleConnectionsPerHost == 0 {
		maxIdleConnectionsPerHost = DefaultMaxIdleConnectionsPerHost
	}

	baseTransportConfig := &transport.BaseTransportConfig{
		MaxIdleConnections:        maxIdleConnections,
		MaxConnectionsPerHost:     maxConnectionsPerHost,
		MaxIdleConnectionsPerHost: maxIdleConnectionsPerHost,
	}

	// logger round tripper config
	loggerTransportConfig := &transport.LoggerTransportConfig{
		LogRequest:                       p.Config.GetBool("modules.http.client.log.request.enabled"),
		LogRequestBody:                   p.Config.GetBool("modules.http.client.log.request.body"),
		LogRequestLevel:                  log.FetchLogLevel(p.Config.GetString("modules.http.client.log.request.level")),
		LogResponse:                      p.Config.GetBool("modules.http.client.log.response.enabled"),
		LogResponseBody:                  p.Config.GetBool("modules.http.client.log.response.body"),
		LogResponseLevel:                 log.FetchLogLevel(p.Config.GetString("modules.http.client.log.response.level")),
		LogResponseLevelFromResponseCode: p.Config.GetBool("modules.http.client.log.response.level_from_response"),
	}

	// round tripper
	var roundTripper http.RoundTripper
	roundTripper = transport.NewLoggerTransportWithConfig(
		transport.NewBaseTransportWithConfig(baseTransportConfig),
		loggerTransportConfig,
	)

	// round tripper tracing extension
	if p.Config.GetBool("modules.http.client.trace.enabled") {
		roundTripper = otelhttp.NewTransport(roundTripper, otelhttp.WithTracerProvider(p.TracerProvider))

		p.Logger.Debug().Msg("http client: enabled tracing")
	}

	// round tripper metrics extension
	if p.Config.GetBool("modules.http.client.metrics.collect.enabled") {
		namespace := p.Config.GetString("modules.http.client.metrics.collect.namespace")
		subsystem := p.Config.GetString("modules.http.client.metrics.collect.subsystem")

		var buckets []float64
		if bucketsConfig := p.Config.GetString("modules.http.client.metrics.buckets"); bucketsConfig != "" {
			for _, s := range strings.Split(strings.ReplaceAll(bucketsConfig, " ", ""), ",") {
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					buckets = append(buckets, f)
				}
			}
		}

		roundTripper = transport.NewMetricsTransportWithConfig(
			roundTripper,
			&transport.MetricsTransportConfig{
				Registry:                  p.MetricsRegistry,
				Namespace:                 Sanitize(namespace),
				Subsystem:                 Sanitize(subsystem),
				Buckets:                   buckets,
				NormalizeRequestPath:      p.Config.GetBool("modules.http.client.metrics.normalize.request_path"),
				NormalizeRequestPathMasks: Flip(p.Config.GetStringMapString("modules.http.client.metrics.normalize.request_path_masks")),
				NormalizeResponseStatus:   p.Config.GetBool("modules.http.client.metrics.normalize.response_status"),
			},
		)

		p.Logger.Debug().Msg("http client: enabled metrics")
	}

	return roundTripper
}

// FxHttpClientParam allows injection of the required dependencies in [NewFxHttpClient].
type FxHttpClientParam struct {
	fx.In
	Factory      httpclient.HttpClientFactory
	RoundTripper http.RoundTripper
	Config       *config.Config
}

// NewFxHttpClient returns a new [http.Client].
func NewFxHttpClient(p FxHttpClientParam) (*http.Client, error) {
	timeout := p.Config.GetInt("modules.http.client.timeout")
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return p.Factory.Create(
		httpclient.WithTimeout(time.Duration(timeout)*time.Second),
		httpclient.WithTransport(p.RoundTripper),
	)
}
