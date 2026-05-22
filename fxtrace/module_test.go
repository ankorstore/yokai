package fxtrace_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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
		fxlog.FxLogModule,
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

// TestModuleStopSwallowsFlushAndShutdownErrors verifies that a failing
// ForceFlush/Shutdown from the underlying span processor does NOT propagate
// out of fx.App.Stop(). OTel flush/shutdown is best-effort by convention; a
// saturated or restarting collector must not turn a graceful pod shutdown into
// a non-zero exit (which Kubernetes interprets as a crashed pod).
func TestModuleStopSwallowsFlushAndShutdownErrors(t *testing.T) {
	// Use a non-test processor type so the OnStop hook calls both ForceFlush
	// AND Shutdown (the test-processor branch skips Shutdown).
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fx.Decorate(factory.NewFailingTracerProviderFactory),
		fx.Invoke(func(oteltrace.TracerProvider) {}),
	).RequireStart()

	// fx.App.Stop must return nil even though both ForceFlush and Shutdown
	// of the underlying tracer provider return an error.
	assert.NoError(t, app.Stop(context.Background()))
}

// TestModuleStopSwallowsFlushErrorOnTestProcessor covers the test-processor
// branch where Shutdown is skipped but ForceFlush still runs. A failing
// ForceFlush must not surface from Stop().
func TestModuleStopSwallowsFlushErrorOnTestProcessor(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("PROCESSOR_TYPE", "noop")
	t.Setenv("SAMPLER_TYPE", "always-on")

	app := fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fxtrace.FxTraceModule,
		fx.Decorate(factory.NewFailingTracerProviderFactory),
		fx.Invoke(func(oteltrace.TracerProvider) {}),
	).RequireStart()

	assert.NoError(t, app.Stop(context.Background()))
}
