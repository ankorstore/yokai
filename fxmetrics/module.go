package fxmetrics

import (
	"github.com/ankorstore/yokai/log"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "metrics"

// FxMetricsModule is the [Fx] metrics module.
//
// [Fx]: https://github.com/uber-go/fx
var FxMetricsModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewDefaultMetricsRegistryFactory,
		NewFxMetricsRegistry,
	),
)

// FxMetricsRegistryParam allows injection of the required dependencies in [NewFxMetricsRegistry].
type FxMetricsRegistryParam struct {
	fx.In
	Factory    MetricsRegistryFactory
	Logger     *log.Logger
	Collectors []prometheus.Collector `group:"metrics-collectors"`
}

// NewFxMetricsRegistry returns a [prometheus.Registry].
func NewFxMetricsRegistry(p FxMetricsRegistryParam) (*prometheus.Registry, error) {
	registry, err := p.Factory.Create()
	if err != nil {
		p.Logger.Error().Err(err).Msg("failed to create metrics registry")

		return nil, err
	}

	for _, collector := range p.Collectors {
		err = registry.Register(collector)
		if err != nil {
			p.Logger.Error().Err(err).Msgf("failed to register metrics collector %+T", collector)

			return nil, err
		} else {
			p.Logger.Debug().Msgf("registered metrics collector %+T", collector)
		}
	}

	return registry, err
}
