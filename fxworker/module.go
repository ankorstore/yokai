package fxworker

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/worker"
	"github.com/prometheus/client_golang/prometheus"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "worker"

// FxWorkerModule is the [Fx] worker module.
//
// [Fx]: https://github.com/uber-go/fx
var FxWorkerModule = fx.Module(
	ModuleName,
	fx.Provide(
		worker.NewDefaultWorkerPoolFactory,
		NewFxWorkerRegistry,
		NewFxWorkerPool,
		fx.Annotate(
			NewFxWorkerModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

// FxWorkerPoolParam allows injection of the required dependencies in [NewFxWorkerPool].
type FxWorkerPoolParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Generator       uuid.UuidGenerator
	TracerProvider  oteltrace.TracerProvider
	Factory         worker.WorkerPoolFactory
	Config          *config.Config
	Registry        *WorkerRegistry
	Logger          *log.Logger
	MetricsRegistry *prometheus.Registry
}

// NewFxWorkerPool returns a new [worker.WorkerPool].
func NewFxWorkerPool(p FxWorkerPoolParam) (*worker.WorkerPool, error) {
	// logger
	logger := log.FromZerolog(p.Logger.ToZerolog().With().Str("module", ModuleName).Logger())

	// tracer provider
	tracerProvider := worker.AnnotateTracerProvider(p.TracerProvider)

	// config
	deferredStartThreshold := p.Config.GetFloat64("modules.worker.defer")
	if deferredStartThreshold <= 0 {
		deferredStartThreshold = worker.DefaultDeferredStartThreshold
	}

	maxExecutionsAttempts := p.Config.GetInt("modules.worker.attempts")
	if maxExecutionsAttempts <= 0 {
		maxExecutionsAttempts = worker.DefaultMaxExecutionsAttempts
	}

	// metrics
	workerMetricsNamespace := p.Config.GetString("modules.worker.metrics.collect.namespace")
	if workerMetricsNamespace == "" {
		workerMetricsNamespace = p.Config.AppName()
	}

	workerMetricsSubsystem := p.Config.GetString("modules.worker.metrics.collect.subsystem")
	if workerMetricsSubsystem == "" {
		workerMetricsSubsystem = ModuleName
	}

	workerMetrics := worker.NewWorkerMetrics(workerMetricsNamespace, workerMetricsSubsystem)

	// pool
	workerPool, err := p.Factory.Create(
		worker.WithGenerator(p.Generator),
		worker.WithMetrics(workerMetrics),
		worker.WithGlobalDeferredStartThreshold(deferredStartThreshold),
		worker.WithGlobalMaxExecutionsAttempts(maxExecutionsAttempts),
	)
	if err != nil {
		logger.Error().Err(err).Msg("worker pool creation error")

		return nil, err
	}

	if p.Config.GetBool("modules.worker.metrics.collect.enabled") {
		err = workerMetrics.Register(p.MetricsRegistry)
		if err != nil {
			logger.Error().Err(err).Msg("worker metrics registration error")

			return nil, err
		}
	}

	// registration
	workersRegistrations, err := p.Registry.ResolveWorkersRegistrations()
	if err != nil {
		p.Logger.Error().Err(err).Msg("worker resolution error")

		return nil, err
	}

	workerPool.Register(workersRegistrations...)

	// context preparation
	workerPoolCtx := logger.WithContext(context.Background())
	workerPoolCtx = context.WithValue(workerPoolCtx, trace.CtxKey{}, tracerProvider)

	// lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Debug().Msg("starting worker pool")

			return workerPool.Start(workerPoolCtx)
		},
		OnStop: func(context.Context) error {
			logger.Debug().Msg("stopping worker pool")

			return workerPool.Stop()
		},
	})

	return workerPool, nil
}
