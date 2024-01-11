package factory

import (
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

type TestTracerProviderFactory struct {
	Params FxTraceTestFactoryParam
}

type FxTraceTestFactoryParam struct {
	fx.In
	Exporter tracetest.TestTraceExporter
}

func NewTestTracerProviderFactory(p FxTraceTestFactoryParam) trace.TracerProviderFactory {
	return &TestTracerProviderFactory{
		Params: p,
	}
}

func (f *TestTracerProviderFactory) Create(options ...trace.TracerProviderOption) (*otelsdktrace.TracerProvider, error) {
	return otelsdktrace.NewTracerProvider(
		otelsdktrace.WithSpanProcessor(otelsdktrace.NewBatchSpanProcessor(f.Params.Exporter)),
		otelsdktrace.WithSampler(trace.NewAlwaysOnSampler()),
	), nil
}
