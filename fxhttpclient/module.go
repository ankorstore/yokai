package fxhttpclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpclient"
	"github.com/ankorstore/yokai/httpclient/transport"
	"github.com/ankorstore/yokai/log"
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
		httpclient.NewDefaultHttpClientFactory,
		NewFxHttpClient,
	),
)

// FxHttpClientParam allows injection of the required dependencies in [NewFxHttpClient].
type FxHttpClientParam struct {
	fx.In
	Factory        httpclient.HttpClientFactory
	TracerProvider trace.TracerProvider
	Config         *config.Config
	Logger         *log.Logger
}

// NewFxHttpClient returns a new [http.Client].
func NewFxHttpClient(p FxHttpClientParam) (*http.Client, error) {
	timeout := p.Config.GetInt("modules.http.client.timeout")
	if timeout == 0 {
		timeout = DefaultTimeout
	}

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

	loggerTransportConfig := &transport.LoggerTransportConfig{
		LogRequest:                       p.Config.GetBool("modules.http.client.log.request.enabled"),
		LogRequestBody:                   p.Config.GetBool("modules.http.client.log.request.body"),
		LogRequestLevel:                  log.FetchLogLevel(p.Config.GetString("modules.http.client.log.request.level")),
		LogResponse:                      p.Config.GetBool("modules.http.client.log.response.enabled"),
		LogResponseBody:                  p.Config.GetBool("modules.http.client.log.response.body"),
		LogResponseLevel:                 log.FetchLogLevel(p.Config.GetString("modules.http.client.log.response.level")),
		LogResponseLevelFromResponseCode: p.Config.GetBool("modules.http.client.log.response.level_from_response"),
	}

	var roundTripper http.RoundTripper
	roundTripper = transport.NewLoggerTransportWithConfig(
		transport.NewBaseTransportWithConfig(baseTransportConfig),
		loggerTransportConfig,
	)

	p.Logger.
		Debug().
		Int("timeout", timeout).
		Str("base transport config", fmt.Sprintf("%+v", baseTransportConfig)).
		Str("logger transport config", fmt.Sprintf("%+v", loggerTransportConfig)).
		Msg("http client: applied configs")

	if p.Config.GetBool("modules.http.client.trace.enabled") {
		roundTripper = otelhttp.NewTransport(roundTripper, otelhttp.WithTracerProvider(p.TracerProvider))

		p.Logger.Debug().Msg("http client: enabled tracing")
	}

	return p.Factory.Create(
		httpclient.WithTimeout(time.Duration(timeout)*time.Second),
		httpclient.WithTransport(roundTripper),
	)
}
