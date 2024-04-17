package fxgrpcserver

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/grpcserver"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/ankorstore/yokai/healthcheck"
	"github.com/ankorstore/yokai/log"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

const (
	ModuleName         = "grpcserver"
	DefaultAddress     = ":50051"
	DefaultBufconnSize = 1024 * 1024
)

// FxGrpcServerModule is the [Fx] grpcserver module.
//
// [Fx]: https://github.com/uber-go/fx
var FxGrpcServerModule = fx.Module(
	ModuleName,
	fx.Provide(
		grpcserver.NewDefaultGrpcServerFactory,
		NewFxGrpcBufconnListener,
		NewFxGrpcServerRegistry,
		NewFxGrpcServer,
		fx.Annotate(
			NewFxGrpcServerModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

// FxGrpcBufconnListenerParam allows injection of the required dependencies in [NewFxGrpcBufconnListener].
type FxGrpcBufconnListenerParam struct {
	fx.In
	Config *config.Config
}

// NewFxGrpcBufconnListener returns a new [bufconn.Listener].
func NewFxGrpcBufconnListener(p FxGrpcBufconnListenerParam) *bufconn.Listener {
	size := p.Config.GetInt("modules.grpc.server.test.bufconn.size")
	if size == 0 {
		size = DefaultBufconnSize
	}

	return grpcservertest.NewBufconnListener(size)
}

// FxGrpcServerParam allows injection of the required dependencies in [NewFxGrpcBufconnListener].
type FxGrpcServerParam struct {
	fx.In
	LifeCycle       fx.Lifecycle
	Factory         grpcserver.GrpcServerFactory
	Generator       uuid.UuidGenerator
	Listener        *bufconn.Listener
	Registry        *GrpcServerRegistry
	Config          *config.Config
	Logger          *log.Logger
	Checker         *healthcheck.Checker
	TracerProvider  trace.TracerProvider
	MetricsRegistry *prometheus.Registry
}

// NewFxGrpcServer returns a new [grpc.Server].
//
//nolint:cyclop
func NewFxGrpcServer(p FxGrpcServerParam) (*grpc.Server, error) {
	// server interceptors
	unaryInterceptors, streamInterceptors := createInterceptors(p)

	for _, unaryInterceptor := range p.Registry.ResolveGrpcServerUnaryInterceptors() {
		unaryInterceptors = append(unaryInterceptors, unaryInterceptor.HandleUnary())
	}

	for _, streamInterceptor := range p.Registry.ResolveGrpcServerStreamInterceptors() {
		streamInterceptors = append(streamInterceptors, streamInterceptor.HandleStream())
	}

	// server options
	grpcServerOptions := append(
		[]grpc.ServerOption{
			grpc.ChainUnaryInterceptor(unaryInterceptors...),
			grpc.ChainStreamInterceptor(streamInterceptors...),
		},
		p.Registry.ResolveGrpcServerOptions()...,
	)

	// server
	grpcServer, err := p.Factory.Create(
		grpcserver.WithServerOptions(grpcServerOptions...),
		grpcserver.WithReflection(p.Config.GetBool("modules.grpc.server.reflection.enabled")),
	)
	if err != nil {
		return nil, err
	}

	// server healthcheck registration
	if p.Config.GetBool("modules.grpc.server.healthcheck.enabled") {
		grpcServer.RegisterService(&grpc_health_v1.Health_ServiceDesc, grpcserver.NewGrpcHealthCheckService(p.Checker))
	}

	// server services registration
	resolvedServices, err := p.Registry.ResolveGrpcServerServices()
	if err != nil {
		return nil, err
	}

	for _, service := range resolvedServices {
		grpcServer.RegisterService(service.Description(), service.Implementation())
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := p.Config.GetString("modules.grpc.server.address")
			if address == "" {
				address = DefaultAddress
			}

			go func() {
				var lis net.Listener
				if p.Config.IsTestEnv() {
					lis = p.Listener
				} else {
					lis, err = net.Listen("tcp", address)
					if err != nil {
						p.Logger.Error().Err(err).Msgf("failed to listen on %s for grpc server", address)
					}
				}

				p.Logger.Info().Msgf("grpc server starting on %s", address)

				if err = grpcServer.Serve(lis); err != nil {
					p.Logger.Error().Err(err).Msg("failed to serve grpc server")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !p.Config.IsTestEnv() {
				grpcServer.GracefulStop()
			}

			return nil
		},
	})

	return grpcServer, nil
}

//nolint:cyclop
func createInterceptors(p FxGrpcServerParam) ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	// panic recovery
	panicRecoveryHandler := grpcserver.NewGrpcPanicRecoveryHandler()

	// interceptors
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(
			recovery.WithRecoveryHandlerContext(panicRecoveryHandler.Handle(p.Config.AppDebug())),
		),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		recovery.StreamServerInterceptor(
			recovery.WithRecoveryHandlerContext(panicRecoveryHandler.Handle(p.Config.AppDebug())),
		),
	}

	// tracer
	if p.Config.GetBool("modules.grpc.server.trace.enabled") {
		var methodFilters []otelgrpc.Filter
		for _, method := range p.Config.GetStringSlice("modules.grpc.server.trace.exclude") {
			methodFilters = append(methodFilters, filters.FullMethodName(method))
		}

		//nolint:staticcheck
		unaryInterceptors = append(
			unaryInterceptors,
			otelgrpc.UnaryServerInterceptor(
				otelgrpc.WithTracerProvider(p.TracerProvider),
				otelgrpc.WithInterceptorFilter(filters.None(methodFilters...)),
			),
		)

		//nolint:staticcheck
		streamInterceptors = append(
			streamInterceptors,
			otelgrpc.StreamServerInterceptor(
				otelgrpc.WithTracerProvider(p.TracerProvider),
				otelgrpc.WithInterceptorFilter(filters.None(methodFilters...)),
			),
		)
	}

	// logger
	loggerInterceptor := grpcserver.
		NewGrpcLoggerInterceptor(p.Generator, log.FromZerolog(p.Logger.ToZerolog().With().Str("system", ModuleName).Logger())).
		Metadata(p.Config.GetStringMapString("modules.grpc.server.log.metadata")).
		Exclude(p.Config.GetStringSlice("modules.grpc.server.log.exclude")...)

	unaryInterceptors = append(unaryInterceptors, loggerInterceptor.UnaryInterceptor())
	streamInterceptors = append(streamInterceptors, loggerInterceptor.StreamInterceptor())

	// metrics
	if p.Config.GetBool("modules.grpc.server.metrics.collect.enabled") {
		namespace := p.Config.GetString("modules.grpc.server.metrics.collect.namespace")
		subsystem := p.Config.GetString("modules.grpc.server.metrics.collect.subsystem")

		var grpcSrvMetricsSubsystemParts []string
		if namespace != "" {
			grpcSrvMetricsSubsystemParts = append(grpcSrvMetricsSubsystemParts, namespace)
		}
		if subsystem != "" {
			grpcSrvMetricsSubsystemParts = append(grpcSrvMetricsSubsystemParts, subsystem)
		}

		grpcSrvMetricsSubsystem := Sanitize(strings.Join(grpcSrvMetricsSubsystemParts, "_"))

		var grpcSrvMetricsBuckets []float64
		if bucketsConfig := p.Config.GetString("modules.grpc.server.metrics.buckets"); bucketsConfig != "" {
			for _, s := range Split(bucketsConfig) {
				f, err := strconv.ParseFloat(s, 64)
				if err == nil {
					grpcSrvMetricsBuckets = append(grpcSrvMetricsBuckets, f)
				}
			}
		}

		if len(grpcSrvMetricsBuckets) == 0 {
			grpcSrvMetricsBuckets = prometheus.DefBuckets
		}

		grpcSrvMetrics := grpcprom.NewServerMetrics(
			grpcprom.WithServerCounterOptions(
				grpcprom.WithSubsystem(grpcSrvMetricsSubsystem),
			),
			grpcprom.WithServerHandlingTimeHistogram(
				grpcprom.WithHistogramSubsystem(grpcSrvMetricsSubsystem),
				grpcprom.WithHistogramBuckets(grpcSrvMetricsBuckets),
			),
		)

		p.MetricsRegistry.MustRegister(grpcSrvMetrics)

		exemplar := func(ctx context.Context) prometheus.Labels {
			if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
				return prometheus.Labels{
					"traceID": span.TraceID().String(),
					"spanID":  span.SpanID().String(),
				}
			}

			return nil
		}

		unaryInterceptors = append(
			unaryInterceptors,
			grpcSrvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplar)),
		)

		streamInterceptors = append(
			streamInterceptors,
			grpcSrvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplar)),
		)
	}

	return unaryInterceptors, streamInterceptors
}
