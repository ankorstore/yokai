package factory

import (
	"context"
	"errors"

	"github.com/ankorstore/yokai/trace"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ErrFailingSpanProcessor is returned by failingSpanProcessor's ForceFlush and Shutdown.
// It mirrors the OTLP "sending_queue is full" failure mode that motivates the
// best-effort shutdown contract in fxtrace.
var ErrFailingSpanProcessor = errors.New("traces export: sending_queue is full")

// failingSpanProcessor is a span processor whose ForceFlush and Shutdown both
// return an error, used to verify fxtrace swallows those errors during OnStop.
type failingSpanProcessor struct{}

func (failingSpanProcessor) OnStart(context.Context, otelsdktrace.ReadWriteSpan) {}
func (failingSpanProcessor) OnEnd(otelsdktrace.ReadOnlySpan)                     {}
func (failingSpanProcessor) ForceFlush(context.Context) error                    { return ErrFailingSpanProcessor }
func (failingSpanProcessor) Shutdown(context.Context) error                      { return ErrFailingSpanProcessor }

// FailingTracerProviderFactory builds a tracer provider whose ForceFlush and
// Shutdown will fail (because the registered span processor fails). It is used
// by tests that exercise the OnStop best-effort path.
type FailingTracerProviderFactory struct{}

func NewFailingTracerProviderFactory() trace.TracerProviderFactory {
	return &FailingTracerProviderFactory{}
}

func (f *FailingTracerProviderFactory) Create(options ...trace.TracerProviderOption) (*otelsdktrace.TracerProvider, error) {
	return otelsdktrace.NewTracerProvider(
		otelsdktrace.WithSpanProcessor(failingSpanProcessor{}),
		otelsdktrace.WithSampler(trace.NewAlwaysOnSampler()),
	), nil
}
