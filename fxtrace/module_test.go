package fxtrace_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxtrace"
	"github.com/ankorstore/yokai/fxtrace/testdata/factory"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleWithTestEnv(t *testing.T) {
	// should fall back on test processor type when APP_ENV=test, no matter given config
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}

func TestModuleSafetyFallbackOnNoopProcessor(t *testing.T) {
	// should fall back on noop processor
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "otlp-grpc")
	t.Setenv("SAMPLER_TYPE", "always-on")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	assert.False(t, exporter.HasSpan("test span", attribute.String("test attribute name", "test attribute value")))
}

func TestModuleWithTestProcessorAndParentBasedAlwaysOnSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "parent-based-always-on")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}

func TestModuleWithTestProcessorAndParentBasedAlwaysOffSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "parent-based-always-off")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	assert.False(t, exporter.HasSpan("test span", attribute.String("test attribute name", "test attribute value")))
}

func TestModuleWithTestProcessorAndParentBasedTraceIdRatioSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "parent-based-trace-id-ratio")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}

func TestModuleWithTestProcessorAndAlwaysOnSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "always-on")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}

func TestModuleWithTestProcessorAndAlwaysOffSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "always-off")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	assert.False(t, exporter.HasSpan("test span", attribute.String("test attribute name", "test attribute value")))
}

func TestModuleWithTestProcessorAndTraceIdRatioSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "test")
	t.Setenv("SAMPLER_TYPE", "trace-id-ratio")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}

func TestModuleWithNoopProcessorAndAlwaysOnSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			defer span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	assert.False(t, exporter.HasSpan("test span", attribute.String("test attribute name", "test attribute value")))
}

func TestModuleWithStdoutProcessorAndAlwaysOffSampler(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "stdout")
	t.Setenv("SAMPLER_TYPE", "always-off")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	assert.False(t, exporter.HasSpan("test span", attribute.String("test attribute name", "test attribute value")))
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")

	var exporter tracetest.TestTraceExporter

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxtrace.FxTraceModule,
		fx.Decorate(factory.NewTestTracerProviderFactory),
		fx.Invoke(func(tracerProvider oteltrace.TracerProvider) {
			_, span := tracerProvider.Tracer("test tracer").Start(
				context.Background(),
				"test span",
				oteltrace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
			)
			span.End()
		}),
		fx.Populate(&exporter),
	).RequireStart().RequireStop()

	tracetest.AssertHasTraceSpan(t, exporter, "test span", attribute.String("test attribute name", "test attribute value"))
}
