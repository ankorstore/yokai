package fxmetrics

import (
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	Config     *config.Config
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

	var registrableCollectors []prometheus.Collector

	if p.Config.GetBool("modules.metrics.collect.build") {
		registrableCollectors = append(registrableCollectors, collectors.NewBuildInfoCollector())
	}

	if p.Config.GetBool("modules.metrics.collect.process") {
		registrableCollectors = append(registrableCollectors, collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	if p.Config.GetBool("modules.metrics.collect.go") {
		registrableCollectors = append(registrableCollectors, collectors.NewGoCollector())
	}

	registrableCollectors = append(registrableCollectors, p.Collectors...)

	for _, collector := range registrableCollectors {
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
