package fxtrace

import (
	"context"
	"fmt"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "trace"

// shutdownCap bounds best-effort flush/shutdown calls so a hanging or
// saturated collector cannot consume the entire pod termination grace period.
const shutdownCap = 5 * time.Second

// FxTraceModule is the [Fx] trace module.
//
// [Fx]: https://github.com/uber-go/fx
var FxTraceModule = fx.Module(
	ModuleName,
	fx.Provide(
		trace.NewDefaultTracerProviderFactory,
		tracetest.NewDefaultTestTraceExporter,
		fx.Annotate(
			NewFxTracerProvider,
			fx.As(new(oteltrace.TracerProvider)),
		),
	),
)

// FxTraceParam allows injection of the required dependencies in [NewFxTracerProvider].
type FxTraceParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Factory   trace.TracerProviderFactory
	Exporter  tracetest.TestTraceExporter
	Config    *config.Config
	Logger    *log.Logger
}

// NewFxTracerProvider returns a [otelsdktrace.TracerProvider].
func NewFxTracerProvider(p FxTraceParam) (*otelsdktrace.TracerProvider, error) {
	ctx := context.Background()

	res, err := createResource(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("cannot create tracer provider resource: %w", err)
	}

	proc, err := createSpanProcessor(ctx, p)
	if err != nil {
		// safety fallback to noop span processor
		proc = trace.NewNoopSpanProcessor()
	}

	samp := createSampler(p)

	tracerProvider, err := p.Factory.Create(
		trace.WithResource(res),
		trace.WithSpanProcessor(proc),
		trace.WithSampler(samp),
	)
	if err != nil {
		return nil, err
	}

	logger := log.FromZerolog(p.Logger.ToZerolog().With().Str("module", ModuleName).Logger())

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// Telemetry flush/shutdown is best-effort by OTel convention: a
			// saturated or restarting collector must not turn a graceful pod
			// shutdown into a non-zero exit. Log and swallow.
			bestEffortStop(ctx, "force flush", tracerProvider.ForceFlush, logger)

			if fetchSpanProcessorType(p) != trace.TestSpanProcessor {
				bestEffortStop(ctx, "shutdown", tracerProvider.Shutdown, logger)
			}

			return nil
		},
	})

	return tracerProvider, nil
}

// bestEffortStop runs a best-effort telemetry operation under a bounded context
// and logs any error without propagating it.
func bestEffortStop(parent context.Context, name string, fn func(context.Context) error, logger *log.Logger) {
	ctx, cancel := boundedContext(parent, shutdownCap)
	defer cancel()

	if err := fn(ctx); err != nil {
		logger.Warn().Err(err).Msgf("tracer provider %s failed (suppressed)", name)
	}
}

// boundedContext derives a child context capped at the smaller of `limit` and
// half of the parent's remaining deadline. The fraction leaves room for any
// subsequent shutdown work so a single hanging exporter cannot consume the
// entire grace period.
func boundedContext(parent context.Context, limit time.Duration) (context.Context, context.CancelFunc) {
	timeout := limit

	if deadline, ok := parent.Deadline(); ok {
		if half := time.Until(deadline) / 2; half > 0 && half < timeout {
			timeout = half
		}
	}

	if timeout <= 0 {
		return context.WithCancel(parent)
	}

	return context.WithTimeout(parent, timeout)
}

func fetchSpanProcessorType(p FxTraceParam) trace.SpanProcessor {
	if p.Config.IsTestEnv() {
		return trace.TestSpanProcessor
	} else {
		return trace.FetchSpanProcessor(p.Config.GetString("modules.trace.processor.type"))
	}
}

func createResource(ctx context.Context, p FxTraceParam) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(p.Config.AppName()),
		),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func createSpanProcessor(ctx context.Context, p FxTraceParam) (otelsdktrace.SpanProcessor, error) {
	switch fetchSpanProcessorType(p) {
	case trace.StdoutSpanProcessor:
		var opts []stdouttrace.Option
		if p.Config.GetBool("modules.trace.processor.options.pretty") {
			opts = append(opts, stdouttrace.WithPrettyPrint())
		}

		return trace.NewStdoutSpanProcessor(opts...), nil
	case trace.TestSpanProcessor:
		return trace.NewTestSpanProcessor(p.Exporter), nil
	case trace.OtlpGrpcSpanProcessor:
		conn, err := trace.NewOtlpGrpcClientConnection(ctx, p.Config.GetString("modules.trace.processor.options.host"))
		if err != nil {
			return nil, err
		}

		return trace.NewOtlpGrpcSpanProcessor(ctx, conn)
	default:
		return trace.NewNoopSpanProcessor(), nil
	}
}

func createSampler(p FxTraceParam) otelsdktrace.Sampler {
	sampler := trace.FetchSampler(p.Config.GetString("modules.trace.sampler.type"))

	switch sampler {
	case trace.ParentBasedAlwaysOffSampler:
		return trace.NewParentBasedAlwaysOffSampler()
	case trace.ParentBasedTraceIdRatioSampler:
		return trace.NewParentBasedTraceIdRatioSampler(p.Config.GetFloat64("modules.trace.sampler.options.ratio"))
	case trace.AlwaysOnSampler:
		return trace.NewAlwaysOnSampler()
	case trace.AlwaysOffSampler:
		return trace.NewAlwaysOffSampler()
	case trace.TraceIdRatioSampler:
		return trace.NewTraceIdRatioSampler(p.Config.GetFloat64("modules.trace.sampler.options.ratio"))
	default:
		return trace.NewParentBasedAlwaysOnSampler()
	}
}
